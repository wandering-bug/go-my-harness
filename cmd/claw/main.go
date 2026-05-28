// cmd/claw/main.go
package main

import (
    "context"
    "log"
    "os"

    "github.com/mambo-wang/go-my-harness/internal/engine"
    "github.com/mambo-wang/go-my-harness/internal/provider"
    "github.com/mambo-wang/go-my-harness/internal/tools"
)

func main() {
    // 确保已设置 MINIMAX_API_KEY, 我的贡献出来给大家用！
    //export MINIMAX_API_KEY="sk-cp-B7wLOi1u7D10BybqBfS50vmPufJ_e88g4arwKLkuDnIH6WpO4MElIO-MCgvf1hOZgyErLd-iYNCNpJqYIoaRdmNL40o_3tXBDK4iKNgZoPr1fQ7W8R7H5WI"
    if os.Getenv("MINIMAX_API_KEY") == "" {
        log.Fatal("请先导出 MINIMAX_API_KEY 环境变量")
    }

    // 1. 获取工作区物理边界
    workDir, _ := os.Getwd()

    // 2. 初始化真实的大脑 (指向智谱 MiniMax-M2.7，使用上一讲的 OpenAI 适配器)
    llmProvider := provider.NewMiniMaxOpenAIProvider("MiniMax-M2.7")

    registry := tools.NewRegistry()
    registry.Register(tools.NewReadFileTool(workDir))
    registry.Register(tools.NewWriteFileTool(workDir))
    registry.Register(tools.NewBashTool(workDir))

    eng := engine.NewAgentEngine(llmProvider, registry, workDir, false)

    prompt := `
请帮我执行以下操作：

1. 用 bash 查看一下我当前电脑的 Go 版本。
2. 帮我写一个简单的 helloworld.go 文件，输出 "Hello, go-my-harness!"。
3. 用 bash 编译并运行这个 go 文件，确认它能正常工作。
`

    err := eng.Run(context.Background(), prompt)
    if err != nil {
        log.Fatalf("引擎运行崩溃: %v", err)
    }
}