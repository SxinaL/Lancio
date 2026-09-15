package entity

// Vocabulary 单词模型
type Vocabulary struct {
	ID               int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Word             string `json:"word" gorm:"uniqueIndex;not null"`
	Phonetic         string `json:"phonetic" gorm:"default:''"`
	Translation      string `json:"translation" gorm:"not null;default:''"`
	ExampleSentence  string `json:"example_sentence" gorm:"default:''"`
	VocabularyBankID int64  `json:"vocabulary_bank_id" gorm:"index;not null;default:1"`
	ReviewCount      int64  `json:"review_count" gorm:"default:0"`             // 累计复习次数
	ReviewStatus     string `json:"review_status" gorm:"default:'unselected'"` // 本轮复习状态：unselected/pending/reviewed
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
	IncrementReviewCount(id int64) error                    // 累计复习次数 +1
	GetVocabularyForReview(count int) ([]Vocabulary, error) // （旧）按复习次数升序抽取
	SetReviewStatus(id int64, status string) error          // 设置单词本轮复习状态
	StartNewReviewRound(count int) ([]Vocabulary, error)    // 开启新一轮：重置reviewed→unselected，按规则选词并标记pending
	GetPendingReviewCount() (int64, error)                  // 获取待复习单词数量
	GetActiveReviewWords() ([]Vocabulary, error)            // 获取本轮活跃单词（pending+reviewed），用于恢复进度
	ResetReviewRound() error                                // 将本轮所有reviewed单词重置为unselected，为下一轮准备
}
