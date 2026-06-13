// provider/config.go
// 从 config.json 配置文件读取应用配置（模型 + 飞书）
package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ModelConfig 单个模型的配置
type ModelConfig struct {
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
	Provider string `json:"provider"` // "openai" 或 "anthropic"
}

// FeishuConfig 飞书机器人配置
type FeishuConfig struct {
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	EncryptKey  string `json:"encrypt_key"`
	VerifyToken string `json:"verify_token"`
}

// AppConfig config.json 顶层结构
type AppConfig struct {
	DefaultModel string                 `json:"default_model"`
	Models       map[string]ModelConfig `json:"models"`
	Feishu       *FeishuConfig          `json:"feishu,omitempty"`
}

// LoadConfig 从指定路径加载 config.json 配置文件
func LoadConfig(configPath string) (*AppConfig, error) {
	// 如果是相对路径，则基于可执行文件所在目录解析
	if !filepath.IsAbs(configPath) {
		// 先尝试基于当前工作目录
		if _, err := os.Stat(configPath); err == nil {
			// 文件存在，直接使用
		} else {
			// 尝试基于可执行文件所在目录
			exe, err := os.Executable()
			if err == nil {
				alt := filepath.Join(filepath.Dir(exe), configPath)
				if _, err := os.Stat(alt); err == nil {
					configPath = alt
				}
			}
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %w", configPath, err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("配置文件解析失败: %w", err)
	}

	if cfg.DefaultModel == "" {
		return nil, fmt.Errorf("配置文件中未指定 default_model")
	}
	if len(cfg.Models) == 0 {
		return nil, fmt.Errorf("配置文件中未定义任何模型")
	}

	return &cfg, nil
}

// GetModelConfig 获取指定模型的配置，如果 name 为空则返回默认模型
func (c *AppConfig) GetModelConfig(name string) (*ModelConfig, string, error) {
	if name == "" {
		name = c.DefaultModel
	}
	mc, ok := c.Models[name]
	if !ok {
		return nil, "", fmt.Errorf("模型 %q 未在配置文件中定义", name)
	}
	return &mc, name, nil
}
