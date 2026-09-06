package method

import (
	"Lancio/database/entity"
)

// GetAllBanks 获取所有词库
func (db *DB) GetAllVocabularyBanks() ([]entity.VocabularyBank, error) {
	var vocabs []entity.VocabularyBank
	err := db.conn.Order("id").Find(&vocabs).Error
	return vocabs, err
}

// GetBankByID 根据ID获取词库
func (db *DB) GetVocabularyBankByID(id int64) (entity.VocabularyBank, error) {
	var v entity.VocabularyBank
	err := db.conn.First(&v, id).Error
	return v, err
}

// CreateVocabularyBank 创建词库
func (db *DB) CreateVocabularyBank(vocabularyBank *entity.VocabularyBank) error {
	return db.conn.Create(vocabularyBank).Error
}

// DeleteVocabularyBank 删除词库
func (db *DB) DeleteVocabularyBankByID(id int64) error {
	return db.conn.Delete(&entity.VocabularyBank{}, id).Error
}