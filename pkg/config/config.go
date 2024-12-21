package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ProxyHosts []string `json:"proxy_hosts"`
	Cookie     string   `json:"cookie"`
}

// 解析 cookie 字符串为 map
func ParseCookies(cookieStr string) map[string]string {
	cookies := make(map[string]string)
	if cookieStr == "" {
		return cookies
	}

	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		cookies[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	
	return cookies
}

var GlobalConfig Config

func LoadConfig(configPath string) error {
	// 如果未指定配置文件路径，使用默认路径
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".pixiv-download", "config.json")
	}

	// 确保配置文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建默认配置
		defaultConfig := Config{
			ProxyHosts: []string{},
			Cookie:    "",
		}
		
		// 确保目录存在
		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return err
		}
		
		// 写入默认配置
		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()
		
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(defaultConfig); err != nil {
			return err
		}
	}

	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&GlobalConfig); err != nil {
		return err
	}

	return nil
} 