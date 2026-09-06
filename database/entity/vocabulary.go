package entity

// Vocabulary 单词模型
type Vocabulary struct {
	ID              int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Word            string `json:"word" gorm:"uniqueIndex;not null"`
	Phonetic        string `json:"phonetic" gorm:"default:''"`
	Translation     string `json:"translation" gorm:"not null;default:''"`
	ExampleSentence string `json:"example_sentence" gorm:"default:''"`
	VocabularyBankID    int64  `json:"vocabulary_bank_id" gorm:"index;not null;default:1"`
}

// VocabularyAPI 单词仓库接口
type VocabularyAPI interface {
	GetAll() ([]Vocabulary, error)
	GetByID(id int64) (Vocabulary, error)
	GetRandom() (Vocabulary, error)
	GetByRange(offset, limit int) ([]Vocabulary, error)
	GetByWord(word string) (Vocabulary, error)
	GetByVocabularyBankID(vocabularyBankID int64) ([]Vocabulary, error)
	DeleteVocabularyByID(id int64) error
	UpdateVocabulary(vocabulary *Vocabulary) error
	CountByVocabularyBankID(vocabularyBankID int64) (int, error)
	CountAll() (int, error)
	AddVocabulary(vocabulary *Vocabulary) error
}
