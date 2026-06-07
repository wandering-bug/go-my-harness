// provider/config.go
// 从 models.json 配置文件读取模型配置
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

// ModelsConfig models.json 顶层结构
type ModelsConfig struct {
	DefaultModel string                 `json:"default_model"`
	Models       map[string]ModelConfig `json:"models"`
}

// LoadModelsConfig 从指定路径加载 models.json 配置文件
func LoadModelsConfig(configPath string) (*ModelsConfig, error) {
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

	var cfg ModelsConfig
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
func (c *ModelsConfig) GetModelConfig(name string) (*ModelConfig, string, error) {
	if name == "" {
		name = c.DefaultModel
	}
	mc, ok := c.Models[name]
	if !ok {
		return nil, "", fmt.Errorf("模型 %q 未在配置文件中定义", name)
	}
	return &mc, name, nil
}
