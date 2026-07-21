package services

import "fmt"

// WordDetail 单词详细信息（扩展语义）
type WordDetail struct {
	Word            string   `json:"word"`
	PhoneticUS      string   `json:"phonetic_us"`
	PhoneticUK      string   `json:"phonetic_uk"`
	Translation     string   `json:"translation"`
	PartOfSpeech    string   `json:"part_of_speech"`
	Definitions     []string `json:"definitions"`
	Examples        []string `json:"examples"`
	Synonyms        []string `json:"synonyms"`
	Antonyms        []string `json:"antonyms"`
}

// DictionaryService 词典服务接口
// 这是一个预留接口，用于后续接入在线词典API（如有道、金山词霸等）
type DictionaryService interface {
	// Lookup 查询单词的详细信息
	Lookup(word string) (*WordDetail, error)
}

// DefaultDictionaryService 默认词典服务实现（本地简易版）
type DefaultDictionaryService struct{}

// NewDefaultDictionaryService 创建默认词典服务
func NewDefaultDictionaryService() *DefaultDictionaryService {
	return &DefaultDictionaryService{}
}

// Lookup 查询单词详细信息（当前返回占位信息）
// TODO: 后续接入在线词典API，返回完整的词义信息
func (s *DefaultDictionaryService) Lookup(word string) (*WordDetail, error) {
	return nil, fmt.Errorf("词典服务尚未配置: 请接入在线词典API以获取'%s'的详细信息", word)
}
