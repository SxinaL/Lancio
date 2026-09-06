package entity

// VocabularyBank 词库模型
type VocabularyBank struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"default:''"`
}

// VocabularyBankAPI 词库API接口
type VocabularyBankAPI interface {
	GetAllVocabularyBanks() ([]VocabularyBank, error)
	GetVocabularyBankByID(id int64) (VocabularyBank, error)		
	CreateVocabularyBank(vocabularyBank *VocabularyBank) error
	DeleteVocabularyBankByID(id int64) error
}
