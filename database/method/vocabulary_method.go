package method

import (
	"fmt"
	"math/rand"

	"Lancio/database/entity"
)

// GetAll 获取所有单词
func (db *DB) GetAll() ([]entity.Vocabulary, error) {
	var words []entity.Vocabulary
	err := db.conn.Order("id").Find(&words).Error
	return words, err
}

// GetByID 根据ID获取单词
func (db *DB) GetByID(id int64) (entity.Vocabulary, error) {
	var w entity.Vocabulary
	err := db.conn.First(&w, id).Error
	return w, err
}

// GetRandom 获取随机单词
func (db *DB) GetRandom() (entity.Vocabulary, error) {
	var count int64
	if err := db.conn.Model(&entity.Vocabulary{}).Count(&count).Error; err != nil {
		return entity.Vocabulary{}, err
	}
	if count == 0 {
		return entity.Vocabulary{}, fmt.Errorf("词库为空")
	}

	offset := rand.Intn(int(count))
	var w entity.Vocabulary
	err := db.conn.Offset(offset).Limit(1).First(&w).Error
	return w, err
}

// GetRandomByVocabularyBankID 从指定词库中获取随机单词
func (db *DB) GetRandomByVocabularyBankID(vocabularyBankID int64) (entity.Vocabulary, error) {
	var count int64
	if err := db.conn.Model(&entity.Vocabulary{}).Where("vocabulary_bank_id = ?", vocabularyBankID).Count(&count).Error; err != nil {
		return entity.Vocabulary{}, err
	}
	if count == 0 {
		return entity.Vocabulary{}, fmt.Errorf("该词库为空")
	}

	offset := rand.Intn(int(count))
	var w entity.Vocabulary
	err := db.conn.Where("vocabulary_bank_id = ?", vocabularyBankID).Offset(offset).Limit(1).First(&w).Error
	return w, err
}

// GetByRange 获取指定范围的单词
func (db *DB) GetByRange(offset, limit int) ([]entity.Vocabulary, error) {
	var words []entity.Vocabulary
	err := db.conn.Order("id").Offset(offset).Limit(limit).Find(&words).Error
	return words, err
}

// Count 获取单词总数
func (db *DB) CountAll() (int, error) {
	var count int64
	err := db.conn.Model(&entity.Vocabulary{}).Count(&count).Error
	return int(count), err
}

// GetByVocabularyBank 根据词库ID获取单词
func (db *DB) GetByVocabularyBankID(vocabularyBankID int64) ([]entity.Vocabulary, error) {
	var words []entity.Vocabulary
	err := db.conn.Where("vocabulary_bank_id = ?", vocabularyBankID).Order("id").Find(&words).Error
	return words, err
}

// DeleteVocabulary 根据ID删除单词
func (db *DB) DeleteVocabulary(id int64) error {
	return db.conn.Delete(&entity.Vocabulary{}, id).Error
}

// UpdateVocabulary 更新单词信息
func (db *DB) UpdateVocabulary(vocabulary *entity.Vocabulary) error {
	return db.conn.Model(&entity.Vocabulary{}).Where("id = ?", vocabulary.ID).Updates(map[string]interface{}{
		"word":             vocabulary.Word,
		"phonetic":         vocabulary.Phonetic,
		"translation":      vocabulary.Translation,
		"example_sentence": vocabulary.ExampleSentence,
		"vocabulary_bank_id": vocabulary.VocabularyBankID,
	}).Error
}

// CountByVocabularyBankID 根据词库ID获取单词数量
func (db *DB) CountByVocabularyBankID(vocabularyBankID int64) (int, error) {
	var count int64
	err := db.conn.Model(&entity.Vocabulary{}).Where("vocabulary_bank_id = ?", vocabularyBankID).Count(&count).Error
	return int(count), err
}

// GetByWord 根据单词获取单词
func (db *DB) GetByWord(word string) (entity.Vocabulary, error) {
	var w entity.Vocabulary
	err := db.conn.Where("word = ?", word).First(&w).Error
	return w, err
}

// AddVocabulary 新增单词
func (db *DB) AddVocabulary(vocabulary *entity.Vocabulary) error {
	return db.conn.Create(vocabulary).Error
}