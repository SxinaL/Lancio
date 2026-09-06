package utils

import (

	"encoding/json"
	"fmt"
	"regexp"

	"strings"

	"github.com/dop251/goja"
)

// MakeLParam 将 x, y 坐标转换为 Windows 消息的 lParam 参数
func ToLParam(x, y int) uintptr {
	// 1. 确保坐标在 16 位范围内 (0-65535)，正确处理负数
	//    例如：-19 & 0xFFFF = 0xFFED
	lx := uint32(x & 0xFFFF)
	ly := uint32(y & 0xFFFF)

	// 2. 将 Y 左移 16 位放到高16位，然后与 X 进行按位或操作
	lParam := (ly << 16) | lx

	// 3. 返回 uintptr 类型，以便直接用于 Windows API 调用
	return uintptr(lParam)
}

// ExtractNuxtData 从有道返回的 HTML 中提取 window.__NUXT__ 数据并返回 JSON 字符串。
//
// 有道词典使用 Nuxt.js SSR，把数据封装在 IIFE 里反混淆：
//
//	window.__NUXT__=(function(a,b,c,...){return {layout:"search",data:[...],...}})("v1","v2",...)
//
// 本函数用栈匹配平衡括号定位 IIFE 边界（处理字符串内的括号和转义），
// 再用 goja 执行 JS，最后 JSON.stringify 返回。
// 能应对参数顺序变化、结构变化、字符串内含特殊字符等情况。
//
// 返回的 JSON 形如：
//
//	{"layout":"search","data":[{"q":"hello","wordData":{...}}],"state":{...}}
func ExtractNuxtData(html []byte) (string, error) {
	const marker = `window.__NUXT__=`
	s := string(html)
	idx := strings.Index(s, marker)
	if idx < 0 {
		return "", fmt.Errorf("HTML 中未找到 window.__NUXT__")
	}

	// IIFE 从 marker 后第一个 '(' 开始
	start := idx + len(marker)
	openIdx := strings.IndexByte(s[start:], '(')
	if openIdx < 0 {
		return "", fmt.Errorf("IIFE 格式异常：缺少左括号")
	}
	openIdx += start

	// 用栈匹配平衡括号，考虑字符串和转义
	depth := 0
	inString := false
	var strCh byte
	escaped := false
	endIdx := -1

	for i := openIdx; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if inString {
			switch ch {
			case '\\':
				escaped = true
			case strCh:
				inString = false
			}
			continue
		}
		switch ch {
		case '"', '\'':
			inString = true
			strCh = ch
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				endIdx = i + 1 // 含右括号
			}
		}
		if endIdx > 0 {
			break
		}
	}
	if endIdx < 0 {
		return "", fmt.Errorf("IIFE 未闭合")
	}

	iife := s[openIdx:endIdx]

	// 用 goja 执行 IIFE 并 JSON.stringify
	vm := goja.New()
	script := "JSON.stringify(" + iife + ")"
	val, err := vm.RunString(script)
	if err != nil {
		return "", fmt.Errorf("执行 IIFE 失败: %w", err)
	}

	jsonStr, ok := val.Export().(string)
	if !ok {
		return "", fmt.Errorf("IIFE 执行结果非字符串: %T", val.Export())
	}
	return jsonStr, nil
}

// ExtractValidFields 从 Nuxt 完整 JSON 中提取有效字段，返回能直接
// json.Unmarshal 到 model.Result 的 JSON 字符串。
//
// 思路：Nuxt 原始结构是 {data:[{q, wordData:{...}}]}，
// model.Result 的 JSON tag 期望扁平结构 {q, word, lang, pronounce, phrs, ...}。
// 本函数把 wordData 下的字段提升到 data[0] 顶层，让 tag 直接对齐。
//
// 字段映射（路径见 model.Result 注释）：
//
//	data[0].q                                  -> q
//	data[0].wordData.phrs.word                 -> word
//	data[0].wordData.lang                     -> lang
//	data[0].wordData.simple.word[0]            -> pronounce
//	data[0].wordData.phrs.phrs                 -> phrs
//	data[0].wordData.ec.word.wfs               -> wfs
//	data[0].wordData.ec.word.trs               -> trs
//	data[0].wordData.blng_sents_part.sentence-pair -> sentence-pair
//
// 容错：任何字段缺失都会跳过，不影响其他字段。
func ExtractValidFields(jsonStr string) (string, error) {
	

	var tmp map[string]interface{}=make(map[string]interface{})
	var root map[string]interface{}

	if err := json.Unmarshal([]byte(jsonStr), &root); err != nil {
		return "", err
	}

	dataArr, ok := root["data"].([]interface{})
	if !ok || len(dataArr) == 0 {
		return "", fmt.Errorf("data 数组为空")
	}
	data0, ok := dataArr[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data 数组第一个元素不是 map")
	}
	//	data[0].q                                  -> q
	tmp["q"] = data0["q"]

	wordData, ok := data0["wordData"].(map[string]interface{})
	if v, ok := wordData["lang"].(string); ok {
		tmp["lang"] = v
	}

	// data[0].wordData.phrs.{word, phrs} -> data[0].{word, phrs}
	if phrs, ok := wordData["phrs"].(map[string]interface{}); ok {
		if v, ok := phrs["word"].(string); ok {
			tmp["word"] = v
		}
		if v, ok := phrs["phrs"].([]interface{}); ok {
			tmp["phrs"] = v
		}
	}

	// data[0].wordData.simple.word[0] -> data[0].pronounce
	if simple, ok := wordData["simple"].(map[string]interface{}); ok {
		if words, ok := simple["word"].([]interface{}); ok && len(words) > 0 {
			tmp["pronounce"] = words[0]
		}
	}

	// data[0].wordData.ec.word.{wfs, trs} -> data[0].{wfs, trs}
	// wfs 原始形式: [{wf:{name, value}}, ...]
	// 反序列化需要: [{name, value}, ...]
	// 所以把每个元素的 "wf" 内层提升到元素顶层，再删掉 "wf"
	if ec, ok := wordData["ec"].(map[string]interface{}); ok {
		if word, ok := ec["word"].(map[string]interface{}); ok {
			if v, ok := word["wfs"].([]interface{}); ok {
				wfs := make([]interface{}, 0, len(v))
				for _, item := range v {
					m, ok := item.(map[string]interface{})
					if !ok {
						continue
					}
					if wf, ok := m["wf"].(map[string]interface{}); ok {
						// 用 wf 的内容替换原元素
						for k, val := range wf {
							m[k] = val
						}
						delete(m, "wf")
					}
					wfs = append(wfs, m)
				}
				tmp["wfs"] = wfs
			}
			if v, ok := word["trs"].([]interface{}); ok {
				tmp["trs"] = v
			}
		}
	}

	// data[0].wordData.blng_sents_part.sentence-pair -> data[0].sentence-pair
	if blng, ok := wordData["blng_sents_part"].(map[string]interface{}); ok {
		if v, ok := blng["sentence-pair"].([]interface{}); ok {
			tmp["sentence-pair"] = v
		}
	}

	// 删掉 wordData 避免多余字段干扰反序列化
	// delete(data0, "wordData")

	tmpByte, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return string(tmpByte), nil
}

// IsLongText 判断是否为长文本。
//
// 规则：短语至多 3 个空格（即 4 个单词），超过的都算句子/长文本。
//
// 注意：返回 true 表示长文本（走百度翻译），false 表示短语（走有道词典）。
func IsLongText(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	spaceCount := strings.Count(text, " ")
	return spaceCount > 3
}

func IsValidWord(s string, language string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z'\- ]*$`, s)
	return matched
}
