package database

// Word 单词模型
type Word struct {
	ID              int64  `json:"id"`
	Word            string `json:"word"`
	Phonetic        string `json:"phonetic"`
	Translation     string `json:"translation"`
	PartOfSpeech    string `json:"part_of_speech"`
	ExampleSentence string `json:"example_sentence"`
	Difficulty      int    `json:"difficulty"` // 1-5
}

// WordRepository 单词仓库接口
type WordRepository interface {
	// GetAll 获取所有单词
	GetAll() ([]Word, error)
	// GetByID 根据ID获取单词
	GetByID(id int64) (Word, error)
	// GetRandom 获取随机单词
	GetRandom() (Word, error)
	// GetByRange 获取指定范围的单词 (分页)
	GetByRange(offset, limit int) ([]Word, error)
	// Count 获取单词总数
	Count() (int, error)
}
