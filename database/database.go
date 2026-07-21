package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB 数据库连接
type DB struct {
	conn *sql.DB
}

// New 创建新的数据库连接
func New(dbPath string) (*DB, error) {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 设置连接参数
	conn.SetMaxOpenConns(1)

	db := &DB{conn: conn}
	if err := db.init(); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	return db, nil
}

// init 初始化数据库表
func (db *DB) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL UNIQUE,
		phonetic TEXT DEFAULT '',
		translation TEXT NOT NULL DEFAULT '',
		part_of_speech TEXT DEFAULT '',
		example_sentence TEXT DEFAULT '',
		difficulty INTEGER DEFAULT 1
	);
	`

	if _, err := db.conn.Exec(schema); err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}

	// 检查是否有数据，如果没有则插入默认词库
	var count int
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM words").Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		return db.seedDefaultWords()
	}

	return nil
}

// seedDefaultWords 插入默认词库
func (db *DB) seedDefaultWords() error {
	words := GetDefaultWords()
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO words (word, phonetic, translation, part_of_speech, example_sentence, difficulty) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, w := range words {
		if _, err := stmt.Exec(w.Word, w.Phonetic, w.Translation, w.PartOfSpeech, w.ExampleSentence, w.Difficulty); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	return db.conn.Close()
}

// GetAll 获取所有单词
func (db *DB) GetAll() ([]Word, error) {
	rows, err := db.conn.Query("SELECT id, word, phonetic, translation, part_of_speech, example_sentence, difficulty FROM words ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Word, &w.Phonetic, &w.Translation, &w.PartOfSpeech, &w.ExampleSentence, &w.Difficulty); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

// GetByID 根据ID获取单词
func (db *DB) GetByID(id int64) (Word, error) {
	var w Word
	err := db.conn.QueryRow(
		"SELECT id, word, phonetic, translation, part_of_speech, example_sentence, difficulty FROM words WHERE id = ?", id,
	).Scan(&w.ID, &w.Word, &w.Phonetic, &w.Translation, &w.PartOfSpeech, &w.ExampleSentence, &w.Difficulty)
	return w, err
}

// GetRandom 获取随机单词
func (db *DB) GetRandom() (Word, error) {
	var count int
	if err := db.conn.QueryRow("SELECT COUNT(*) FROM words").Scan(&count); err != nil {
		return Word{}, err
	}
	if count == 0 {
		return Word{}, fmt.Errorf("词库为空")
	}
	offset := rand.Intn(count)
	var w Word
	err := db.conn.QueryRow(
		"SELECT id, word, phonetic, translation, part_of_speech, example_sentence, difficulty FROM words LIMIT 1 OFFSET ?", offset,
	).Scan(&w.ID, &w.Word, &w.Phonetic, &w.Translation, &w.PartOfSpeech, &w.ExampleSentence, &w.Difficulty)
	return w, err
}

// GetByRange 获取指定范围的单词
func (db *DB) GetByRange(offset, limit int) ([]Word, error) {
	rows, err := db.conn.Query(
		"SELECT id, word, phonetic, translation, part_of_speech, example_sentence, difficulty FROM words ORDER BY id LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Word, &w.Phonetic, &w.Translation, &w.PartOfSpeech, &w.ExampleSentence, &w.Difficulty); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

// Count 获取单词总数
func (db *DB) Count() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	return count, err
}

// Ensure DB implements WordRepository
var _ WordRepository = (*DB)(nil)
