package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/pipeline"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// getConfigPath returns the path to the configs directory
func getConfigPath(t *testing.T) string {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// If we're in test directory, go up one level
	if filepath.Base(wd) == "test" {
		return filepath.Join(filepath.Dir(wd), "configs", "pipelines")
	}

	// Otherwise assume we're in project root
	return filepath.Join(wd, "configs", "pipelines")
}

// TestPipelineParser 测试流水线解析器
func TestPipelineParser(t *testing.T) {
	parser := pipeline.NewParser(getConfigPath(t))

	t.Run("LoadCodeSimplePipeline", func(t *testing.T) {
		p, err := parser.LoadPipeline("code_simple")
		if err != nil {
			t.Fatalf("Failed to load pipeline: %v", err)
		}

		if p.Name != "code_simple" {
			t.Errorf("Expected name 'code_simple', got '%s'", p.Name)
		}

		if p.WorkflowType != "CODE_SIMPLE" {
			t.Errorf("Expected workflow_type 'CODE_SIMPLE', got '%s'", p.WorkflowType)
		}

		if len(p.Nodes) == 0 {
			t.Error("Expected at least one node")
		}
	})

	t.Run("ValidatePipeline", func(t *testing.T) {
		p, err := parser.LoadPipeline("code_simple")
		if err != nil {
			t.Fatalf("Failed to load pipeline: %v", err)
		}

		// Validation should pass
		if err := parser.ValidatePipeline(p); err != nil {
			t.Errorf("Pipeline validation failed: %v", err)
		}
	})

	t.Run("CalculateExecutionOrder", func(t *testing.T) {
		p, err := parser.LoadPipeline("code_simple")
		if err != nil {
			t.Fatalf("Failed to load pipeline: %v", err)
		}

		order, err := parser.CalculateExecutionOrder(p)
		if err != nil {
			t.Fatalf("Failed to calculate execution order: %v", err)
		}

		if len(order) != len(p.Nodes) {
			t.Errorf("Expected %d nodes, got %d", len(p.Nodes), len(order))
		}

		// Check that discovery comes before clarification
		var discoveryIdx, clarificationIdx int
		for i, nodeOrder := range order {
			if nodeOrder.Node.Ref == "discovery" {
				discoveryIdx = i
			}
			if nodeOrder.Node.Ref == "clarification" {
				clarificationIdx = i
			}
		}

		if discoveryIdx >= clarificationIdx {
			t.Error("Expected discovery to come before clarification")
		}
	})

	t.Run("LoadAllPipelines", func(t *testing.T) {
		pipelines, err := parser.LoadAllPipelines()
		if err != nil {
			t.Fatalf("Failed to load all pipelines: %v", err)
		}

		if len(pipelines) == 0 {
			t.Error("Expected at least one pipeline")
		}

		if _, exists := pipelines["code_simple"]; !exists {
			t.Error("Expected code_simple pipeline to exist")
		}
	})

	t.Run("ListPipelines", func(t *testing.T) {
		infos, err := parser.ListPipelines()
		if err != nil {
			t.Fatalf("Failed to list pipelines: %v", err)
		}

		if len(infos) == 0 {
			t.Error("Expected at least one pipeline info")
		}

		// Check that info contains required fields
		for _, info := range infos {
			if info.Name == "" {
				t.Error("Pipeline name is required")
			}
			if info.DisplayName == "" {
				t.Error("Pipeline display_name is required")
			}
		}
	})

	t.Run("GetPipelineByWorkflowType", func(t *testing.T) {
		p, err := parser.GetPipelineByWorkflowType("CODE_SIMPLE")
		if err != nil {
			t.Fatalf("Failed to get pipeline by workflow type: %v", err)
		}

		if p.WorkflowType != "CODE_SIMPLE" {
			t.Errorf("Expected workflow_type 'CODE_SIMPLE', got '%s'", p.WorkflowType)
		}
	})

	t.Run("DetectCircularDependencies", func(t *testing.T) {
		// Create a pipeline with circular dependency
		p := &types.Pipeline{
			Name:         "circular",
			DisplayName:  "Circular Pipeline",
			Version:      "0.1.0",
			WorkflowType: "CODE_SIMPLE",
			Nodes: []types.PipelineNode{
				{
					Ref:         "a",
					DisplayName: "Node A",
					Agent:       "investigator",
					Model:       "haiku",
					Timeout:     60,
					Required:    true,
					DependsOn:   []string{"b"},
				},
				{
					Ref:         "b",
					DisplayName: "Node B",
					Agent:       "investigator",
					Model:       "haiku",
					Timeout:     60,
					Required:    true,
					DependsOn:   []string{"a"},
				},
			},
		}

		err := parser.CheckCircularDependencies(p)
		if err == nil {
			t.Error("Expected error for circular dependency")
		}
	})
}

// TestNodeValidation 测试节点验证
func TestNodeValidation(t *testing.T) {
	parser := pipeline.NewParser("../../configs/pipelines")

	t.Run("ValidNode", func(t *testing.T) {
		node := &types.PipelineNode{
			Ref:         "test",
			DisplayName: "Test Node",
			Agent:       "investigator",
			Model:       "haiku",
			Timeout:     60,
			Required:    true,
		}

		err := parser.ValidateNode(node, 0)
		if err != nil {
			t.Errorf("Expected valid node, got error: %v", err)
		}
	})

	t.Run("MissingRef", func(t *testing.T) {
		node := &types.PipelineNode{
			DisplayName: "Test Node",
			Agent:       "investigator",
			Model:       "haiku",
			Timeout:     60,
			Required:    true,
		}

		err := parser.ValidateNode(node, 0)
		if err == nil {
			t.Error("Expected error for missing ref")
		}
	})

	t.Run("InvalidTimeout", func(t *testing.T) {
		node := &types.PipelineNode{
			Ref:         "test",
			DisplayName: "Test Node",
			Agent:       "investigator",
			Model:       "haiku",
			Timeout:     -1,
			Required:    true,
		}

		err := parser.ValidateNode(node, 0)
		if err == nil {
			t.Error("Expected error for invalid timeout")
		}
	})
}
