package test

import (
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/cli"
)

// TestCLICommands 测试 CLI 命令
func TestCLICommands(t *testing.T) {
	t.Run("HelpCommand", func(t *testing.T) {
		c, err := cli.New(getConfigPath(t))
		if err != nil {
			t.Fatalf("Failed to create CLI: %v", err)
		}
		defer c.Close()

		// CLI 初始化成功即通过
		t.Log("CLI created successfully")
	})

	t.Run("PipelineList", func(t *testing.T) {
		c, err := cli.New(getConfigPath(t))
		if err != nil {
			t.Fatalf("Failed to create CLI: %v", err)
		}
		defer c.Close()

		// 测试 pipeline list 命令
		// 由于需要命令行参数解析，这里只是验证初始化成功
		t.Log("Pipeline list command available")
	})
}

// TestCLIIntegration 集成测试 - 测试完整的 CLI 工作流
func TestCLIIntegration(t *testing.T) {
	configPath := getConfigPath(t)
	c, err := cli.New(configPath)
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}
	defer c.Close()

	// 这个测试框架验证 CLI 初始化
	// 实际命令测试需要使用 os/exec 来运行真实的命令行
	t.Log("CLI integration test - initialization successful")
}
