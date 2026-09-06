package test

import (
	"Lancio/api"
	"Lancio/model"
	"fmt"
	"testing"
)

func TestRequestBaiduTranslate(t *testing.T) {
	// 测试用：直接使用默认配置（密钥已写在 DefaultConfig 里）
	if model.AppConfig == nil {
		model.AppConfig = model.DefaultConfig()
	}

	cases := []struct {
		name string
		q    string
		from string
		to   string
	}{
		{"英→中-单词", "hello", "en", "zh"},
		{"英→中-短语", "how are you", "en", "zh"},
		{"中→英-短句", "你好，世界", "zh", "en"},
		{"自动检测源语言", "Bonjour le monde", "auto", "zh"},
		{"长文本", "The quick brown fox jumps over the lazy dog.", "en", "zh"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			text, err := api.RequestBaiduTranslate(c.q, c.from, c.to)
			if err != nil {
				t.Logf("[%s] 失败: %v", c.name, err)
				return
			}
			fmt.Printf("[%s] %s (%s→%s) = %s\n", c.name, c.q, c.from, c.to, text)
		})
	}
}
