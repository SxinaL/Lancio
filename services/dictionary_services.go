package services

import (
	"Lancio/api"
	"Lancio/database/entity"
	"Lancio/model"
	"Lancio/utils"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// FetchLocal 先查本地缓存。
//
// 命中：返回反序列化后的 *model.Result（Found=true）
// 未命中或读取失败：返回非 nil error，调用方应继续走在线查询
//
// 长文本不走 SQLite 缓存。
func FetchLocal(query, language string) (*model.Result, error) {
	if model.AppInstance == nil || model.AppInstance.Db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	db := model.AppInstance.Db

	searchResult, err := db.Get(query, language)
	if err != nil {
		return nil, fmt.Errorf("词典查询发送错误: %w", err)
	}
	if searchResult.Detail == nil {
		return nil, fmt.Errorf("未找到当前查询结果: %s", query)
	}

	// zlib 解压
	data, err := zlib.NewReader(bytes.NewReader(searchResult.Detail))
	if err != nil {
		return nil, fmt.Errorf("zlib解压失败: %w", err)
	}
	defer data.Close()

	jsonbyte, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("读取单词失败: %w", err)
	}
	if len(jsonbyte) == 0 {
		return nil, fmt.Errorf("单词详细为空")
	}

	// 反序列化
	result := &model.Result{}
	if err := json.Unmarshal(jsonbyte, result); err != nil {
		return nil, fmt.Errorf("unmarshal失败: %w", err)
	}
	return result, nil
}

// FetchOnline 在线查询，返回解析后的 *model.Result
func FetchOnline(query, language string) (*model.Result, error) {
	jsonStr, err := api.RequestYoudao(query, language)
	if err != nil {
		return nil, err
	}

	result := &model.Result{}
	if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
		return nil, fmt.Errorf("反序列化失败: %w", err)
	}
	// 有时候存在多音标的情况，取第一个
	result.Pronounce.Ukphone = "/" + strings.Split(result.Pronounce.Ukphone, "; ")[0] + "/"
	result.Pronounce.Usphone = "/" + strings.Split(result.Pronounce.Usphone, "; ")[0] + "/"
	return result, nil
}

// LookUpLongText 长文本走百度翻译 API（from=auto, to=zh），返回译文
func LookUpLongText(query string) (string, error) {
	return api.RequestBaiduTranslate(query, "auto", "zh")
}

// LookupWord 短文本查询：内存缓存 -> SQLite -> 有道在线（含回写）
func LookupWord(query string, language string) (*entity.Vocabulary, error) {
	if query == "" {
		return nil, fmt.Errorf("空查询")
	}
	query = strings.TrimSpace(strings.ToLower(query))

	// 1. 内存缓存
	if r := utils.GetResultCache(query, language); r != nil {
		return utils.ToVocabulary(r), nil
	}

	// 2. SQLite 缓存
	if result, err := FetchLocal(query, language); err == nil {
		// 命中：回填内存缓存并返回
		utils.SetResultCache(query, language, *result)
		return utils.ToVocabulary(result), nil
	}

	// 3. 在线查询
	result, err := FetchOnline(query, language)
	if err != nil {
		return utils.ToVocabulary(&model.Result{}), fmt.Errorf("查询失败: %w", err)
	}
	// 无翻译的时候，暂无结果，不写入数据库
	if len(result.Translations) == 0 {
		return utils.ToVocabulary(result), nil
	}
	// 4. 回写 SQLite 缓存（失败不影响返回结果）
	go utils.UpdateDictionary(query, language, result)

	// 5. 回填内存缓存
	utils.SetResultCache(query, language, *result)
	return utils.ToVocabulary(result), nil
}
