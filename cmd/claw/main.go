// cmd/claw/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/mambo-wang/go-my-harness/internal/engine"
    "github.com/mambo-wang/go-my-harness/internal/provider"
    "github.com/mambo-wang/go-my-harness/internal/tools"
)

func main() {
    // 1. 获取工作区物理边界
    workDir, _ := os.Getwd()

    // 2. 从 models.json 配置文件读取模型配置
    configPath := filepath.Join(workDir, "models.json")
    cfg, err := provider.LoadModelsConfig(configPath)
    if err != nil {
        log.Fatalf("加载配置文件失败: %v", err)
    }

    mc, modelName, err := cfg.GetModelConfig("")
    if err != nil {
        log.Fatalf("获取模型配置失败: %v", err)
    }
    fmt.Printf("[Config] 使用模型: %s (provider=%s)\n", modelName, mc.Provider)

    // 3. 根据 provider 类型创建 LLM Provider
    var llmProvider provider.LLMProvider
    switch mc.Provider {
    case "openai":
        llmProvider = provider.NewOpenAIProvider(mc.APIKey, mc.BaseURL, modelName)
    case "anthropic":
        llmProvider = provider.NewAnthropicProvider(mc.APIKey, mc.BaseURL, modelName)
    default:
        log.Fatalf("不支持的 provider 类型: %s", mc.Provider)
    }

    registry := tools.NewRegistry()
    registry.Register(tools.NewReadFileTool(workDir))
    registry.Register(tools.NewWriteFileTool(workDir))
    registry.Register(tools.NewBashTool(workDir))
    registry.Register(tools.NewEditFileTool(workDir))

    eng := engine.NewAgentEngine(llmProvider, registry, workDir, false)

    prompt := `
我当前目录下有一个 server.go 文件。
请帮我把里面 "TODO: 增加鉴权逻辑" 下面的那个 if 语句，整个替换为：

if user == nil {
    fmt.Println("Forbidden!")
    return
}
`
    err = eng.Run(context.Background(), prompt)
    if err != nil {
        log.Fatalf("引擎运行崩溃: %v", err)
    }
}