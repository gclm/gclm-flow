package test

import (
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TestWorkflowParser 测试流水线解析器
func TestWorkflowParser(t *testing.T) {
	parser := workflow.NewParser(getConfigPath(t))

	t.Run("LoadCodeSimpleWorkflow", func(t *testing.T) {
		p, err := parser.LoadWorkflow("code_simple")
		if err != nil {
			t.Fatalf("Failed to load pipeline: %v", err)
		}

		if p.Name != "code_simple" {
			t.Errorf("Expected name 'code_simple', got '%s'", p.Name)
		}

		if p.WorkflowType != "fix" {
			t.Errorf("Expected workflow_type 'fix', got '%s'", p.WorkflowType)
		}

		if len(p.Nodes) == 0 {
			t.Error("Expected at least one node")
		}
	})

	t.Run("ValidateWorkflow", func(t *testing.T) {
		p, err := parser.LoadWorkflow("code_simple")
		if err != nil {
			t.Fatalf("Failed to load pipeline: %v", err)
		}

		// Validation should pass
		if err := parser.ValidateWorkflow(p); err != nil {
			t.Errorf("Workflow validation failed: %v", err)
		}
	})

	t.Run("CalculateExecutionOrder", func(t *testing.T) {
		p, err := parser.LoadWorkflow("code_simple")
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

	t.Run("LoadAllWorkflows", func(t *testing.T) {
		pipelines, err := parser.LoadAllWorkflows()
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

	t.Run("ListWorkflows", func(t *testing.T) {
		infos, err := parser.ListWorkflows()
		if err != nil {
			t.Fatalf("Failed to list pipelines: %v", err)
		}

		if len(infos) == 0 {
			t.Error("Expected at least one pipeline info")
		}

		// Check that info contains required fields
		for _, info := range infos {
			if info.Name == "" {
				t.Error("Workflow name is required")
			}
			if info.DisplayName == "" {
				t.Error("Workflow display_name is required")
			}
		}
	})

	t.Run("GetWorkflowByType", func(t *testing.T) {
		p, err := parser.GetWorkflowByType("fix")
		if err != nil {
			t.Fatalf("Failed to get pipeline by workflow type: %v", err)
		}

		if p.WorkflowType != "fix" {
			t.Errorf("Expected workflow_type 'fix', got '%s'", p.WorkflowType)
		}
	})

	t.Run("DetectCircularDependencies", func(t *testing.T) {
		// Create a pipeline with circular dependency
		p := &types.Workflow{
			Name:         "circular",
			DisplayName:  "Circular Workflow",
			Version:      "0.1.0",
			WorkflowType: "fix",
			Nodes: []types.WorkflowNode{
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
	parser := workflow.NewParser("../../workflows")

	t.Run("ValidNode", func(t *testing.T) {
		node := &types.WorkflowNode{
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
		node := &types.WorkflowNode{
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
		node := &types.WorkflowNode{
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
