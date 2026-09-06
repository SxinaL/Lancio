package model

import (
	"sync"
	"time"
)

// Result 查询结果
type Result struct {
	Query     string `json:"q"`    // 用户输入的查询文本  data[0].q
	Word      string `json:"word"` // 单词 data[0].wordData.phrs.word
	Language  string `json:"lang"` // 语言 data[0].wordData.lang
	Pronounce struct {
		Usphone  string `json:"usphone"`  // 发音：nation -> phonetic，如 "美" -> "/æˈbændən/"
		Usspeech string `json:"usspeech"` // 发音链接
		Ukphone  string `json:"ukphone"`  // 发音：nation -> phonetic，如 "英" -> "/æˈbændən/"
		Ukspeech string `json:"ukspeech"` // 发音链接
	} `json:"pronounce"` // 发音信息 data[0].wordData.simple.word[0]
	Phrases []struct { // 词组列表 data[0].wordData.phrs.phrs
		Headword    string `json:"headword"`    // 词组头词
		Translation string `json:"translation"` // 词组翻译
	} `json:"phrs"`
	WordForms []struct { // 词形变体列表 data[0].wordData.ec.word.wfs
		Name  string `json:"name"`  // 词形变体
		Value string `json:"value"` // 词性
	} `json:"wfs"`
	Translations []struct { // 释义列表 data[0].wordData.ec.word.trs
		Pos  string `json:"pos"`  // 词性
		Tran string `json:"tran"` // 释义
	} `json:"trs"`
	SentencePairs []struct { // 例句对列表 data[0].wordData.blng_sents_part.sentence-pair
		Sentence            string `json:"sentence"`             // 例句
		SentenceTranslation string `json:"sentence-translation"` // 例句翻译
	} `json:"sentence-pair"`
}

// CacheEntry 内存缓存条目
type CacheEntry struct {
	Result  Result
	Expires time.Time
}

// ResultCache 内存级查询结果缓存，避免反复查 SQLite
//
// key 格式: "query|lang"
var ResultCache = struct {
	Mu   sync.RWMutex
	Data map[string]CacheEntry
	TTL  time.Duration
}{
	Data: make(map[string]CacheEntry),
	TTL:  30 * time.Minute, // 默认 TTL，可按需调整
}
