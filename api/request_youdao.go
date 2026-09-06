package api

import (

	"Lancio/utils"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// userAgents 一些常见桌面浏览器 UA，规避有道反爬
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
}

var uaRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	// ydCliNormal 有道词典 HTTP 客户端，用于正常请求
	ydCliNormal = &http.Client{Timeout: 5 * time.Second}
	// ydCliHTTPS 有道词典 HTTPS 客户端，用于降级尝试
	ydCliHTTPS = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   5 * time.Second,
	}
	// youdaoBaseURL = "http://dict.youdao.com/w"
	youdaoBaseURL = "https://dict.youdao.com/result"
)

// GetRandomUA 随机返回一个 UA
func GetRandomUA() string {
	return userAgents[uaRand.Intn(len(userAgents))]
}

// requestYoudao 向有道词典发送 GET 请求,返回响应体
func RequestYoudao(query string, language string) (string, error) {
	query = strings.ReplaceAll(query, " ", "%20")
	url := fmt.Sprintf("%s?word=%s&lang=%s", strings.TrimRight(youdaoBaseURL, "/"), query, language)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "dict.youdao.com")
	req.Header.Set("User-Agent", GetRandomUA())

	resp, err := ydCliNormal.Do(req)
	if err != nil {
		// 降级尝试 HTTPS 客户端
		resp, err = ydCliHTTPS.Do(req)
		if err != nil {
			return "", fmt.Errorf("请求有道失败: %w", err)
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("有道返回 %s", resp.Status)
	}
	if len(body) == 0 {
		return "", fmt.Errorf("有道返回空响应")
	}
	// 调用方可用 ExtractNuxtData(body) 从 HTML 中提取 Nuxt SSR 数据
	jsonStr, err := utils.ExtractNuxtData(body)
	if err != nil {
		return "", fmt.Errorf("提取 Nuxt SSR 数据失败: %w", err)
	}
	// 提取有效字段
	result, err := utils.ExtractValidFields(jsonStr)
	return result, err
}


