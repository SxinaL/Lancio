package method

import (
	"fmt"
	"os"
	"path/filepath"
	"Lancio/database/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB 数据库连接
type DB struct {
	conn *gorm.DB
}

// New 创建新的数据库连接
func New(dbPath string) (*DB, error) {
	// 创建数据库目录（如果不存在）
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}
	// 打开数据库连接
	conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	// 初始化数据库表
	db := &DB{conn: conn}
	if err := db.init(); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	return db, nil
}

// init 初始化数据库表
func (db *DB) init() error {
	// 构建 VocabularyBank 表
	if err := db.conn.AutoMigrate(&entity.VocabularyBank{}); err != nil {
		return fmt.Errorf("构建 VocabularyBank 表失败: %w", err)
	}
	// 构建 Vocabulary 表
	if err := db.conn.AutoMigrate(&entity.Vocabulary{}); err != nil {
		return fmt.Errorf("构建 Vocabulary 表失败: %w", err)
	}

	// 词典查询缓存表：en/ch 两张表结构相同，复用 DictionaryCache
	if err := db.conn.Table("en").AutoMigrate(&entity.Dictionary{}); err != nil {
		return fmt.Errorf("构建 en 缓存表失败: %w", err)
	}
	if err := db.conn.Table("ch").AutoMigrate(&entity.Dictionary{}); err != nil {
		return fmt.Errorf("构建 ch 缓存表失败: %w", err)
	}

	// // 开发使用后期删除默认词库和单词
	// var vocabCount int64
	// if err := db.conn.Model(&entity.VocabularyBank{}).Count(&vocabCount).Error; err != nil {
	// 	return err
	// }

	// if vocabCount == 0 {
	// 	if err := db.seedDefaultVocabulary(); err != nil {
	// 		return err
	// 	}
	// }

	// var wordCount int64
	// if err := db.conn.Model(&entity.Vocabulary{}).Count(&wordCount).Error; err != nil {
	// 	return err
	// }

	// if wordCount == 0 {
	// 	return db.seedDefaultWords()
	// }

	return nil
}

// seedDefaultVocabulary 插入默认词库
func (db *DB) seedDefaultVocabulary() error {
	defaultVocab := &entity.VocabularyBank{Name: "默认词库", Description: "系统默认词库"}
	return db.conn.Create(defaultVocab).Error
}

// seedDefaultWords 插入默认单词
func (db *DB) seedDefaultWords() error {
	words := GetDefaultWords()
	return db.conn.CreateInBatches(words, 100).Error
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

var _ entity.VocabularyBankAPI = (*DB)(nil)
var _ entity.VocabularyBankAPI = (*DB)(nil)
