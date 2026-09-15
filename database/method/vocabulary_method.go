package method

import (
	"fmt"
	"math/rand"

	"Lancio/database/entity"

	"gorm.io/gorm"
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
		"word":               vocabulary.Word,
		"phonetic":           vocabulary.Phonetic,
		"translation":        vocabulary.Translation,
		"example_sentence":   vocabulary.ExampleSentence,
		"vocabulary_bank_id": vocabulary.VocabularyBankID,
		"review_count":       vocabulary.ReviewCount,
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

// IncrementReviewCount 增加单词复习次数 +1
func (db *DB) IncrementReviewCount(id int64) error {
	return db.conn.Model(&entity.Vocabulary{}).Where("id = ?", id).
		UpdateColumn("review_count", gorm.Expr("review_count + 1")).Error
}

// GetVocabularyForReview 获取复习用单词列表：按复习次数升序抽取 count 个
func (db *DB) GetVocabularyForReview(count int) ([]entity.Vocabulary, error) {
	var words []entity.Vocabulary
	err := db.conn.Order("review_count ASC").Limit(count).Find(&words).Error
	return words, err
}

// SetReviewStatus 设置单词复习状态
func (db *DB) SetReviewStatus(id int64, status string) error {
	return db.conn.Model(&entity.Vocabulary{}).Where("id = ?", id).
		Update("review_status", status).Error
}

// StartNewReviewRound 开启新一轮复习

// 状态流转：
// reviewed  → unselected （重置上一轮已完成的）
// pending   → 保留      （上一轮中断的继续复习）
// unselected → pending   （按review_count升序补足到count个）
func (db *DB) StartNewReviewRound(count int) ([]entity.Vocabulary, error) {
	// 第一步：重置上一轮已复习的单词
	db.conn.Model(&entity.Vocabulary{}).Where("review_status = ?", "reviewed").
		Update("review_status", "unselected")

	// 第二步：获取上一轮中断遗留的待复习单词
	var pending []entity.Vocabulary
	db.conn.Where("review_status = ?", "pending").Order("review_count ASC").Find(&pending)

	// 第三步：不足时从复习次数最少的未选中单词中补足
	need := count - len(pending)
	var fresh []entity.Vocabulary
	if need > 0 {
		db.conn.Where("review_status = ?", "unselected").
			Order("review_count ASC").Limit(need).Find(&fresh)
	}

	// 第四步：合并并批量标记为待复习
	all := append(pending, fresh...)
	ids := make([]int64, len(all))
	for i := range all {
		ids[i] = all[i].ID
		all[i].ReviewStatus = "pending"
	}
	if len(ids) > 0 {
		db.conn.Model(&entity.Vocabulary{}).Where("id IN ?", ids).
			Update("review_status", "pending")
	}

	return all, nil
}

// GetPendingReviewCount 获取待复习单词数量
func (db *DB) GetPendingReviewCount() (int64, error) {
	var count int64
	err := db.conn.Model(&entity.Vocabulary{}).Where("review_status = ?", "pending").Count(&count).Error
	return count, err
}

// ResetReviewRound 将本轮所有reviewed单词重置为unselected，用于一轮复习全部完成后清理状态
func (db *DB) ResetReviewRound() error {
	return db.conn.Model(&entity.Vocabulary{}).
		Where("review_status = ?", "reviewed").
		Update("review_status", "unselected").Error
}

// GetActiveReviewWords 获取当前轮次中所有活跃的复习单词（pending + reviewed），用于恢复进度
func (db *DB) GetActiveReviewWords() ([]entity.Vocabulary, error) {
	var words []entity.Vocabulary
	err := db.conn.Where("review_status IN ?", []string{"pending", "reviewed"}).
		Order("review_status ASC, review_count ASC").Find(&words).Error
	return words, err
}
