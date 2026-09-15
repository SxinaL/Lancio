package handler

import (
	"Lancio/api"
	"Lancio/model"
	"Lancio/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

type APIConfigKey struct {
	Appid string `json:"appid"`
	Key   string `json:"key"`
}

func (c *ConfigHandler) GetAPIConfig() APIConfigKey {
	if model.AppConfig == nil {
		utils.LoadConfig()
	}
	return APIConfigKey{
		Appid: model.AppConfig.APIConfigKey.LANCIO_APPID,
		Key:   model.AppConfig.APIConfigKey.LANCIO_API_KEY,
	}
}

func (c *ConfigHandler) SaveAPIConfig(appid, key string) error {
	utils.LoadConfig()
	model.AppConfig.APIConfigKey.LANCIO_APPID = appid
	model.AppConfig.APIConfigKey.LANCIO_API_KEY = key
	return utils.SaveConfig()
}

func (c *ConfigHandler) TestAPIConnection(appid, key string) error {
	return api.TestBaiduConnection(appid, key)
}

func (c *ConfigHandler) EnsureAPIConfigured() bool {
	if model.AppConfig == nil {
		utils.LoadConfig()
	}
	if model.AppConfig.APIConfigKey.LANCIO_APPID != "" && model.AppConfig.APIConfigKey.LANCIO_API_KEY != "" {
		return true
	}
	if model.AppInstance != nil && model.AppInstance.Ctx != nil {
		runtime.MessageDialog(model.AppInstance.Ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "配置提醒",
			Message: "检测到百度翻译 API 密钥尚未配置，请前往设置页面完成配置。",
		})
	}
	return false
}

func (c *ConfigHandler) GetReviewWordsPerSession() int {
	if model.AppConfig == nil {
		utils.LoadConfig()
	}
	if model.AppConfig.ReviewWordsPerSession <= 0 {
		return 10
	}
	return model.AppConfig.ReviewWordsPerSession
}

func (c *ConfigHandler) SetReviewWordsPerSession(count int) error {
	utils.LoadConfig()
	model.AppConfig.ReviewWordsPerSession = count
	return utils.SaveConfig()
}
