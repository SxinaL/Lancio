package method

import (
	"errors"

	"Lancio/database/entity"
)



// Get 查询词典（未命中返回 gorm.ErrRecordNotFound）
func (db *DB) Get(query string, language string) (entity.Dictionary, error) {
	var c entity.Dictionary
	err := db.conn.Table(language).Where("query = ?", query).First(&c).Error
	return c, err
}

// Upsert 插入或更新词典
func (db *DB) Upsert(c *entity.Dictionary, language string) error {
	if c.Query == "" {
		return errors.New("query 不能为空")
	}
	// SQLite 用 ON CONFLICT(query) DO REPLACE
	return db.conn.Table(language).Save(c).Error
}

// Delete 删除词典单词（不存在不报错）
func (db *DB) Delete(query string, language string) error {
	return db.conn.Table(language).Where("query = ?", query).Delete(&entity.Dictionary{}).Error
}

// Count 统计词典总条数
func (db *DB) Count(language string) (int, error) {
	var n int64
	err := db.conn.Table(language).Count(&n).Error
	return int(n), err
}

// 接口实现断言
var _ entity.DictionaryAPI = (*DB)(nil)
