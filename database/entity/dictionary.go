package entity

import "time"

// DictionaryCache 词典查询缓存（en/ch 两张表结构一致，复用同一个结构体）
type Dictionary struct {
	Query      string    `json:"query" gorm:"primaryKey;type:text;not null"`
	Detail     []byte    `json:"detail" gorm:"type:blob;not null"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;not null"`
}

// DictionaryCacheAPI 词典缓存仓库接口
// isEN=true 操作 en 表（英文查询），false 操作 ch 表（中文查询）
type DictionaryAPI interface {
	Get(query string, language string) (Dictionary, error)
	Upsert(dictionary *Dictionary, language string) error
	Delete(query string, language string) error
	Count(language string) (int, error)
}
