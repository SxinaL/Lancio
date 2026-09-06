package model

// WindowFloatConfig 浮窗窗口配置结构
type WindowFloatConfig struct {
	WindowMode               string  `json:"windowMode"`
	WindowOpacity            float64 `json:"windowOpacity"`
	WindowPositionX          int     `json:"windowPositionX"`
	WindowPositionY          int     `json:"windowPositionY"`
	WindowSizeW              int     `json:"windowSizeW"`
	WindowSizeH              int     `json:"windowSizeH"`
	WindowSizeMinW           int     `json:"windowSizeMinW"`
	WindowSizeMinH           int     `json:"windowSizeMinH"`
	SelectedVocabularyBankID int64   `json:"selectedVocabularyBankID"`
	IsPinned                 bool    `json:"isPinned"`
	Launch                   bool    `json:"launch"`
}

// API 配置结构密钥
type APIConfigKey struct {
	LANCIO_APPID   string `json:"appid"`
	LANCIO_API_KEY string `json:"key"`
}

// Config 应用配置结构
type Config struct {
	WindowFloatConfigs map[string]*WindowFloatConfig `json:"windowFloatConfig"`
	APIConfigKey       APIConfigKey                  `json:"apiConfigKey"`
	StartHidden        bool                          `json:"startHidden"`
}

// defaultConfig 返回默认配置（autoStart 默认开启）
func DefaultConfig() *Config {
	return &Config{
		WindowFloatConfigs: map[string]*WindowFloatConfig{
			"LancioFloat": {
				WindowMode:               "float",
				WindowOpacity:            0.4,
				WindowPositionX:          0,
				WindowPositionY:          0,
				WindowSizeW:              280,
				WindowSizeH:              300,
				WindowSizeMinW:           256,
				WindowSizeMinH:           256,
				SelectedVocabularyBankID: 0,
				IsPinned:                 false,
				Launch:                   true,
			},
		},
		APIConfigKey: APIConfigKey{
			LANCIO_APPID:   "",
			LANCIO_API_KEY: "",
		},
		StartHidden: true,
	}
}

var AppConfig *Config
