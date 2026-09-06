package handler

import (

	"Lancio/database/entity"
	"Lancio/services"
)

// DictionaryService 在线词典查询与朗读服务（暴露给前端）
type DictionaryHandler struct{}

func NewDictionaryHandler() *DictionaryHandler {
	return &DictionaryHandler{}
}

// Lookup 在线查询单词或文本，结果填充到 *entity.Vocabulary。
//
//   - 英文单词：返回音标、释义、柯林斯例句
//   - 中->英 句子/长文本：走有道机器翻译
//   - 未能找到：返回 ErrNotFound
func (d *DictionaryHandler) LookupWord(query string, language string) (*entity.Vocabulary, error) {
	return services.LookupWord(query, language)
}

