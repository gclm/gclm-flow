package workflow

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// getWorkflowsPathForBench 找到工作流目录
func getWorkflowsPathForBench() string {
	wd, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		testPath := filepath.Join(wd, "workflows")
		if info, err := os.Stat(testPath); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(testPath)
			if len(entries) > 0 {
				return testPath
			}
		}
		wd = filepath.Dir(wd)
		if filepath.Base(wd) == "gclm-engine" {
			return filepath.Join(wd, "workflows")
		}
	}
	return "workflows"
}

// BenchmarkNewParser 基准测试：创建解析器
func BenchmarkNewParser(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewParser(workflowsDir)
	}
}

// BenchmarkLoadWorkflow 基准测试：加载工作流
func BenchmarkLoadWorkflow(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.LoadWorkflow("code_simple")
	}
}

// BenchmarkLoadAllWorkflows 基准测试：加载所有工作流
func BenchmarkLoadAllWorkflows(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(workflowsDir)
		_, _ = parser.LoadAllWorkflows()
	}
}

// BenchmarkValidateWorkflow 基准测试：验证工作流
func BenchmarkValidateWorkflow(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)
	workflow, _ := parser.LoadWorkflow("code_simple")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.ValidateWorkflow(workflow)
	}
}

// BenchmarkCalculateExecutionOrder 基准测试：计算执行顺序
func BenchmarkCalculateExecutionOrder(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)
	workflow, _ := parser.LoadWorkflow("code_complex")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.CalculateExecutionOrder(workflow)
	}
}

// BenchmarkGetWorkflowByType 基准测试：按类型获取工作流
func BenchmarkGetWorkflowByType(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.GetWorkflowByType("fix")
	}
}

// BenchmarkListWorkflows 基准测试：列出工作流
func BenchmarkListWorkflows(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ListWorkflows()
	}
}

// BenchmarkValidateNode 基准测试：验证节点
func BenchmarkValidateNode(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)
	workflow, _ := parser.LoadWorkflow("code_simple")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range workflow.Nodes {
			_ = parser.ValidateNode(&workflow.Nodes[j], j)
		}
	}
}

// BenchmarkCheckCircularDependencies 基准测试：检查循环依赖
func BenchmarkCheckCircularDependencies(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	parser := NewParser(workflowsDir)
	workflow, _ := parser.LoadWorkflow("code_complex")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.CheckCircularDependencies(workflow)
	}
}

// BenchmarkComplexWorkflow 基准测试：复杂工作流操作
func BenchmarkComplexWorkflow(b *testing.B) {
	// 创建一个复杂的工作流进行测试
	workflow := &types.Workflow{
		Name:         "benchmark",
		DisplayName:  "Benchmark Workflow",
		Version:      "1.0.0",
		WorkflowType: "feat",
		Nodes:        make([]types.WorkflowNode, 20),
	}

	// 填充20个节点，每个有多个依赖
	for i := 0; i < 20; i++ {
		workflow.Nodes[i] = types.WorkflowNode{
			Ref:         "node-" + string(rune('a'+i)),
			DisplayName: "Node " + string(rune('a'+i)),
			Agent:       "investigator",
			Model:       "haiku",
			Timeout:     60,
			Required:    true,
		}
		// 添加依赖
		if i > 0 {
			workflow.Nodes[i].DependsOn = []string{"node-" + string(rune('a'+i-1))}
		}
		if i > 2 {
			workflow.Nodes[i].DependsOn = append(workflow.Nodes[i].DependsOn, "node-"+string(rune('a'+i-2)))
		}
	}

	parser := NewParser("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.CalculateExecutionOrder(workflow)
	}
}

// BenchmarkParallelWorkflowLoading 基准测试：并行加载工作流
func BenchmarkParallelWorkflowLoading(b *testing.B) {
	workflowsDir := getWorkflowsPathForBench()
	workflows := []string{"code_simple", "code_complex", "document"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			parser := NewParser(workflowsDir)
			for _, name := range workflows {
				_, _ = parser.LoadWorkflow(name)
			}
		}
	})
}
