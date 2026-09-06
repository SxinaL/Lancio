package handler

import (
	"Lancio/database/entity"
	"Lancio/model"
	"fmt"
)

type VocabularyHandler struct{}

// NewVocabularyHandler 创建单词处理程序
func NewVocabularyHandler() *VocabularyHandler {
	return &VocabularyHandler{}
}

// GetWordByID 根据ID获取单词
func (vocHandler *VocabularyHandler) GetWordByID(id int64) (entity.Vocabulary, error) {
	if model.AppInstance.Db == nil {
		return entity.Vocabulary{}, fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.GetByID(id)
}

// GetWordsByRange 获取指定范围的单词
func (vocHandler *VocabularyHandler) GetWordsByRange(offset, limit int) ([]entity.Vocabulary, error) {
	if model.AppInstance.Db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.GetByRange(offset, limit)
}

// GetWordCount 获取单词总数
func (vocHandler *VocabularyHandler) GetWordCount() (int, error) {
	if model.AppInstance.Db == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.CountAll()
}

// GetWordsByVocabularyBankID 根据词库ID获取单词列表
func (vocHandler *VocabularyHandler) GetWordsByVocabularyBankID(vocabularyID int64) ([]entity.Vocabulary, error) {
	if model.AppInstance.Db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.GetByVocabularyBankID(vocabularyID)
}

// DeleteWord 根据ID删除单词
func (vocHandler *VocabularyHandler) DeleteVocabulary(id int64) error {
	if model.AppInstance.Db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.DeleteVocabulary(id)
}

// UpdateVocabulary 更新词库信息
func (vocHandler *VocabularyHandler) UpdateVocabulary(w entity.Vocabulary) error {
	if model.AppInstance.Db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.UpdateVocabulary(&w)
}

// AddVocabulary 新增单词
func (vocHandler *VocabularyHandler) AddVocabulary(w entity.Vocabulary) error {
	if model.AppInstance.Db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return model.AppInstance.Db.AddVocabulary(&w)
}

// GetStarStatus 获取单词是否收藏
func (vocHandler *VocabularyHandler) GetStarStatus(word string) (entity.Vocabulary, error) {
	if model.AppInstance.Db == nil {
		return entity.Vocabulary{}, fmt.Errorf("数据库未初始化")
	}
	w, err := model.AppInstance.Db.GetByWord(word)
	return w, err
}
