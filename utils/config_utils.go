package utils

import (
	"Lancio/model"

	"encoding/json"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveConfig 将当前配置写入 config.json
func SaveConfig() error {
	fmt.Println(model.AppConfig.WindowFloatConfigs["LancioFloat"])
	config_path := GetFilePath("config.json")
	data, err := json.MarshalIndent(model.AppConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := WriteToFile(data, config_path); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	fmt.Printf("[Config] 已保存配置文件: %s \n", config_path)
	return nil
}

// GetFloatConfig 获取指定窗口的浮窗配置，键不存在时自动补默认项，避免 nil 解引用
func GetFloatConfig(name string) *model.WindowFloatConfig {
	if model.AppConfig == nil {
		model.AppConfig = model.DefaultConfig()
	}
	if cfg, ok := model.AppConfig.WindowFloatConfigs[name]; ok && cfg != nil {
		return cfg
	}
	cfg := &model.WindowFloatConfig{}
	model.AppConfig.WindowFloatConfigs[name] = cfg
	return cfg
}

// LoadConfig 读取 config.json，读取失败时使用默认配置并创建文件
func LoadConfig() error {
	// appConfigOnce.Do(func() {
	model.AppConfig = model.DefaultConfig()

	config_path := GetFilePath("config.json")
	data, err := os.ReadFile(config_path)
	if err != nil {
		fmt.Printf("[Config] 未找到 config.json，使用默认配置: %v\n", err)
		model.AppConfig = model.DefaultConfig()
		data, err := json.MarshalIndent(model.AppConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化配置失败: %w", err)
		}
		if err := WriteToFile(data, config_path); err != nil {
			return fmt.Errorf("写入配置文件失败: %w", err)
		}

		// 根据默认配置创建 config.json
		fmt.Printf("[Config] 已创建默认配置文件: %s\n", config_path)
		return nil
	}
	// 解析失败也用默认配置覆盖写回，避免下次再失败,这里修改&model.AppConfig
	if err := json.Unmarshal(data, &model.AppConfig); err != nil {
		fmt.Printf("[Config] 解析配置文件失败 %s: %v\n", config_path, err)
		model.AppConfig = model.DefaultConfig()
		data, err := json.MarshalIndent(model.AppConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化配置失败: %w", err)
		}
		if err := WriteToFile(data, config_path); err != nil {
			return fmt.Errorf("重置配置文件失败: %w", err)
		}
		fmt.Printf("[Config] 已重置为默认配置文件: %s\n", config_path)

		return nil
	}
	fmt.Printf("[Config] 已加载配置文件: %s \n", config_path)

	return nil
}

// SaveFloatWindowPositionAndSize 保存浮窗位置和尺寸
func SaveFloatWindowPositionAndSize() {
	if model.AppInstance == nil {
		return
	}
	if GetFloatConfig(model.WindowInstance.WindowName).IsPinned {
		return
	}
	X, Y := runtime.WindowGetPosition(model.AppInstance.Ctx)
	W, H := runtime.WindowGetSize(model.AppInstance.Ctx)
	fmt.Printf("SetFloatWindowPosition: X=%d, Y=%d\n", X, Y)
	LoadConfig()
	cfg := GetFloatConfig(model.WindowInstance.WindowName)
	cfg.WindowPositionX = X
	cfg.WindowPositionY = Y
	cfg.WindowSizeW = W
	cfg.WindowSizeH = H
	SaveConfig()
}
