// internal/engine/reporter.go
package engine

import (
	"context"
	"fmt"
	"strings"
)

// Reporter 定义了 Agent 引擎向外界输出信息的规范。
// 这使得引擎可以无缝切换终端 (CLI)、飞书、钉钉甚至 WebUI 等不同的展现层。
type Reporter interface {
	// OnThinking 当模型开始进行慢思考 (Reasoning) 时调用
	OnThinking(ctx context.Context)

	// OnToolCall 当模型决定并发调用工具时调用
	OnToolCall(ctx context.Context, toolName string, args string)

	// OnToolResult 当工具在底层执行完毕并返回结果时调用
	OnToolResult(ctx context.Context, toolName string, result string, isError bool)

	// OnMessage 当模型宣告任务完成，向用户输出最终纯文本回答时调用
	OnMessage(ctx context.Context, content string)
}

// TerminalReporter 终端模式下的输出报告器，直接打印到控制台
type TerminalReporter struct{}

// NewTerminalReporter 创建一个终端 Reporter
func NewTerminalReporter() *TerminalReporter {
	return &TerminalReporter{}
}

func (r *TerminalReporter) OnThinking(ctx context.Context) {
	fmt.Println("🤔 (模型正在思考...)")
}

func (r *TerminalReporter) OnToolCall(ctx context.Context, toolName string, args string) {
	fmt.Printf("\n🛠️  [工具调用] %s(%s)\n", toolName, args)
}

func (r *TerminalReporter) OnToolResult(ctx context.Context, toolName string, result string, isError bool) {
	prefix := "✅"
	if isError {
		prefix = "❌"
	}
	// 截断过长结果
	if len(result) > 500 {
		result = result[:500] + "...(已截断)"
	}
	fmt.Printf("%s [%s] %s\n", prefix, toolName, strings.TrimSpace(result))
}

func (r *TerminalReporter) OnMessage(ctx context.Context, content string) {
	fmt.Printf("\n📝 %s\n", content)
}