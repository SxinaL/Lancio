package utils

import (
	"Lancio/database/entity"
	"Lancio/model"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"strings"
	"time"
)

// UpdateDictionary 把 Result 序列化为 JSON + zlib 压缩后写入 SQLite
// 用于在线查询成功后异步回写本地缓存
func UpdateDictionary(query, language string, r *model.Result) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(jsonBytes); err != nil {
		w.Close()
		return
	}
	w.Close()

	db := model.AppInstance.Db
	if db == nil {
		return
	}
	_ = db.Upsert(&entity.Dictionary{
		Query:      query,
		Detail:     buf.Bytes(),
		UpdateTime: time.Now(),
	}, language)
}

// ToVocabulary 将 dictionary.Result 转换为 entity.Vocabulary
func ToVocabulary(result *model.Result) *entity.Vocabulary {
	// 单词：优先结构化字段 Word，其次查询文本
	word := result.Query

	vocabulary := &entity.Vocabulary{
		Word: word,
	}

	// 音标：优先美音，其次英音
	if result.Pronounce.Usphone != "" {
		vocabulary.Phonetic = result.Pronounce.Usphone
	} else if result.Pronounce.Ukphone != "" {
		vocabulary.Phonetic = result.Pronounce.Ukphone
	}

	// 翻译：合并所有 Translations（按 "词性 + 释义" 每行一条）
	if len(result.Translations) > 0 {
		parts := make([]string, 0, len(result.Translations))
		for _, tr := range result.Translations {
			// 格式：[pos] tran 或直接 tran（pos 为空时）
			tran := strings.Split(tr.Tran, "；")
			num  := min(2, len(tran))
			tr.Tran = strings.Join(tran[0:num], "；")
			if tr.Pos != "" {
				parts = append(parts, tr.Pos+" "+tr.Tran)
			} else {
				parts = append(parts, tr.Tran)
			}
		}
		vocabulary.Translation = strings.Join(parts, "\n")
	}

	// 例句：取第一条例句对
	if len(result.SentencePairs) > 0 {
		sp := result.SentencePairs[0]
		if sp.Sentence != "" && sp.SentenceTranslation != "" {
			vocabulary.ExampleSentence = sp.Sentence + "\n" + sp.SentenceTranslation
		} else if sp.Sentence != "" {
			vocabulary.ExampleSentence = sp.Sentence
		}
	}

	return vocabulary
}
