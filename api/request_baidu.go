package api

import (

	"Lancio/model"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 百度翻译通用文本翻译 API
// 文档: https://fanyi-api.baidu.com/product/113
//
// 请求方式: POST 或 GET（推荐 POST，q 较长时 GET 会被 URL 长度限制）
// 接口: https://fanyi-api.baidu.com/api/trans/vip/translate
//
// 必需参数:
//   q        待翻译文本，UTF-8 编码
//   from     源语言（可设为 auto）
//   to       目标语言（不可为 auto）
//   appid    应用 ID
//   salt     随机数（防重放）
//   sign     签名 = MD5(appid + q + salt + 密钥)，32 位小写

var (
	baiduEndpoint = "https://fanyi-api.baidu.com/api/trans/vip/translate"

	baiduRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	baiduHTTP = &http.Client{Timeout: 10 * time.Second}
)

// BaiduTransResult 单条翻译结果
type BaiduTransResult struct {
	Src string `json:"src"` // 原文
	Dst string `json:"dst"` // 译文
}

// BaiduTransResponse 百度翻译 API 响应
//
// 成功: {"from":"en","to":"zh","trans_result":[{"src":"hello","dst":"你好"}]}
// 失败: {"error_code":"54001","error_msg":"Invalid Sign"}
type BaiduTransResponse struct {
	From        string             `json:"from"`         // 实际源语言
	To          string             `json:"to"`           // 实际目标语言
	TransResult []BaiduTransResult `json:"trans_result"` // 翻译结果
	ErrorCode   string             `json:"error_code"`   // 错误码（成功时为空）
	ErrorMsg    string             `json:"error_msg"`    // 错误信息（成功时为空）
}

// baiduSign 计算签名: MD5(appid + q + salt + key)，返回 32 位小写 hex
func baiduSign(appID, secretKey, q, salt string) string {
	raw := appID + q + salt + secretKey
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// baiduSalt 生成随机 salt（数字字符串）
func baiduSalt() string {
	return fmt.Sprintf("%d", baiduRand.Int63())
}

// RequestBaiduTranslate 调用百度翻译通用文本翻译 API
//
// 入参:
//
//	q       待翻译文本（UTF-8，无需预先编码）
//	from    源语言代码（如 "en"、"zh"、"auto"）
//	to      目标语言代码（如 "zh"、"en"）
//
// 返回:
//
//	译文拼接（多条结果以 \n 连接），error 为失败原因
func RequestBaiduTranslate(q, from, to string) (string, error) {
	if model.AppConfig == nil {
		return "", fmt.Errorf("配置未初始化: model.AppConfig 为 nil")
	}
	appID := model.AppConfig.APIConfigKey.LANCIO_APPID
	secretKey := model.AppConfig.APIConfigKey.LANCIO_API_KEY
	if appID == "" || secretKey == "" {
		return "", fmt.Errorf("百度翻译密钥未配置: 请在配置文件中设置 APIConfigKey.LANCIO_APPID 和 LANCIO_API_KEY")
	}

	salt := baiduSalt()
	sign := baiduSign(appID, secretKey, q, salt)

	form := url.Values{}
	form.Set("q", q)
	form.Set("from", from)
	form.Set("to", to)
	form.Set("appid", appID)
	form.Set("salt", salt)
	form.Set("sign", sign)

	req, err := http.NewRequest("POST", baiduEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := baiduHTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求百度翻译失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result BaiduTransResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, body=%s", err, string(body))
	}

	if result.ErrorCode != "" {
		return "", fmt.Errorf("百度翻译错误 %s: %s",
			result.ErrorCode, mapErrCode(result.ErrorCode, result.ErrorMsg))
	}

	if len(result.TransResult) == 0 {
		return "", nil
	}

	// 多条结果以 \n 连接，保持原文换行结构
	parts := make([]string, 0, len(result.TransResult))
	for _, r := range result.TransResult {
		parts = append(parts, r.Dst)
	}
	return strings.Join(parts, "\n"), nil
}

// mapErrCode 把常见错误码转成可读说明
//
// 官方错误码参考: https://fanyi-api.baidu.com/doc/24
func mapErrCode(code, msg string) string {
	descriptions := map[string]string{
		"52000": "成功",
		"52001": "请求超时",
		"52002": "系统错误",
		"52003": "未授权用户（appid 错误或服务未开通）",
		"54000": "必填参数为空",
		"54001": "签名错误（检查 appid + key + salt + q 拼接顺序）",
		"54003": "访问频率受限（QPS 超限）",
		"54004": "余额不足",
		"58000": "客户端 IP 非法（不在白名单）",
		"58001": "语言不支持",
		"58002": "服务当前已关闭",
		"90107": "未开通此服务",
	}
	if desc, ok := descriptions[code]; ok {
		return fmt.Sprintf("%s（%s）", desc, msg)
	}
	return msg
}
