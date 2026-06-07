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

    // 实例化引擎，开启 EnableThinking = true (开启慢思考，促使模型一次性统筹规划)
    eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

    // 下发一个需要收集多源信息的任务
    prompt := `
    我当前目录下有 a.txt, b.txt, c.txt 三个文件。
    为了节省时间，请你同时一次性读取这三个文件，并将它们的内容综合起来，告诉我它们分别记录了什么领域的信息。
    `

    err = eng.Run(context.Background(), prompt)
    if err != nil {
        log.Fatalf("引擎运行崩溃: %v", err)
    }
}