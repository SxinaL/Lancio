package handler

import (
	"Lancio/database/entity"
	"Lancio/model"
	"Lancio/utils"
	"fmt"
)

type VocabularyBankHandler struct {}

// NewBankHandler 创建词库处理程序
func NewBankHandler() *VocabularyBankHandler {
	return &VocabularyBankHandler{}
}

// GetAllVocabularyBanks 获取所有词库
func (vb *VocabularyBankHandler) GetAllVocabularyBanks() ([]entity.VocabularyBank, error) {
	db := model.AppInstance.Db
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return db.GetAllVocabularyBanks()
}

// CreateVocabularyBank 创建词库
func (vb *VocabularyBankHandler) CreateVocabularyBank(name, description string) (entity.VocabularyBank, error) {
	db := model.AppInstance.Db
	if db == nil {
		return entity.VocabularyBank{}, fmt.Errorf("数据库未初始化")
	}
	v := &entity.VocabularyBank{Name: name, Description: description}
	if err := db.CreateVocabularyBank(v); err != nil {
		return entity.VocabularyBank{}, err
	}
	return *v, nil
}

// DeleteVocabularyBank 删除词库
func (vb *VocabularyBankHandler) DeleteVocabularyBank(id int64) error {
	db := model.AppInstance.Db
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.DeleteVocabularyBankByID(id)
}

// GetVocabularyCountByVocabularyBankID 获取指定词库的单词数量
func (vb *VocabularyBankHandler) GetVocabularyCountByVocabularyBankID(vocabularyBankID int64) (int, error) {
	db := model.AppInstance.Db
	if db == nil {
		return 0, fmt.Errorf("数据库未初始化")
	}
	return db.CountByVocabularyBankID(vocabularyBankID)
}

// GetRandomVocabulary 获取随机单词：优先从配置选中的词库取，未指定时从全部词库取
func (vb *VocabularyBankHandler) GetRandomVocabulary() (entity.Vocabulary, error) {
	db := model.AppInstance.Db
	if db == nil {
		return entity.Vocabulary{}, fmt.Errorf("数据库未初始化")
	}
	// 读取配置中选中的词库 ID
	if model.AppConfig != nil {
		if cfg := model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName]; cfg != nil && cfg.SelectedVocabularyBankID > 0 {
			return db.GetRandomByVocabularyBankID(cfg.SelectedVocabularyBankID)
		}
	}
	return db.GetRandom()
}

// GetSelectedVocabularyBankID 获取当前配置中选中的加载词库 ID
func (vb *VocabularyBankHandler) GetSelectedVocabularyBankID() (int64, error) {
	if model.AppConfig != nil {
		if cfg := model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName]; cfg != nil {
			return cfg.SelectedVocabularyBankID, nil
		}
	}
	return 0, nil
}

// SetSelectedVocabularyBankID 设置当前加载的词库并持久化到配置
func (vb *VocabularyBankHandler) SetSelectedVocabularyBankID(id int64) error {
	db := model.AppInstance.Db
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	// 校验词库是否存在（id>0 时）
	if id > 0 {
		if _, err := db.GetVocabularyBankByID(id); err != nil {
			return fmt.Errorf("词库不存在: %w", err)
		}
	}

	utils.LoadConfig()
	cfg := model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName]
	if cfg == nil {
		cfg = &model.WindowFloatConfig{}
		model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName] = cfg
	}
	cfg.SelectedVocabularyBankID = id
	return utils.SaveConfig()
}
