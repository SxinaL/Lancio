package main

import (
	"fmt"

	"DesktopVoc/database"
	"DesktopVoc/services"
)

// GetRandomWord 获取一个随机单词
func (a *App) GetRandomWord() (database.Word, error) {
	if a.db == nil {
		return database.Word{}, fmt.Errorf("数据库未初始化")
	}
	return a.db.GetRandom()
}

// GetWordByID 根据ID获取单词
func (a *App) GetWordByID(id int64) (database.Word, error) {
	if a.db == nil {
		return database.Word{}, fmt.Errorf("数据库未初始化")
	}
	return a.db.GetByID(id)
}

// GetWordsByRange 获取指定范围的单词
func (a *App) GetWordsByRange(offset, limit int) ([]database.Word, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return a.db.GetByRange(offset, limit)
}

// GetWordCount 获取单词总数
func (a *App) GetWordCount() (int, error) {
	if a.db == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return a.db.Count()
}

// LookupWordDetail 查询单词详细信息
// 预留接口 - 后续接入在线词典后可以返回完整的词义、例句、同义词等
func (a *App) LookupWordDetail(word string) (*services.WordDetail, error) {
	if a.dict == nil {
		return nil, fmt.Errorf("词典服务未初始化")
	}
	return a.dict.Lookup(word)
}
