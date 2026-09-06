package utils

import (

	"os"
	"path/filepath"


)


// writeConfigToFile 将配置以 JSON 格式写入指定路径
func WriteToFile(data []byte, config_path string) error {
	if err := os.WriteFile(config_path, data, 0644); err != nil {
		return err
	}
	return nil
}




// configPath 返回文件路径（优先当前工作目录，其次可执行文件同目录）
func GetFilePath(fileName string) string {
	if cwd, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(cwd, fileName)); err == nil {
			return filepath.Join(cwd, fileName)
		}
	}

	// 3. fallback：可执行文件同目录
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), fileName)
	}

	p, _ := os.Getwd()
	return filepath.Join(p, fileName)
}