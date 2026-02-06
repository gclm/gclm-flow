package test

import (
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TestWorkflowParser 测试工作流解析器
func TestWorkflowParser(t *testing.T) {
	parser := workflow.NewParser(getWorkflowsPath(t))

	t.Run("LoadFixWorkflow", func(t *testing.T) {
		wf, err := parser.LoadWorkflow("fix")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		if wf.Name != "fix" {
			t.Errorf("Expected name 'fix', got '%s'", wf.Name)
		}

		if wf.WorkflowType != "fix" {
			t.Errorf("Expected workflow_type 'fix', got '%s'", wf.WorkflowType)
		}

		if len(wf.Nodes) == 0 {
			t.Error("Expected at least one node")
		}
	})

	t.Run("ValidateWorkflow", func(t *testing.T) {
		wf, err := parser.LoadWorkflow("fix")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		// Validation should pass
		if err := parser.ValidateWorkflow(wf); err != nil {
			t.Errorf("Workflow validation failed: %v", err)
		}
	})

	t.Run("CalculateExecutionOrder", func(t *testing.T) {
		wf, err := parser.LoadWorkflow("fix")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		order, err := parser.CalculateExecutionOrder(wf)
		if err != nil {
			t.Fatalf("Failed to calculate execution order: %v", err)
		}

		if len(order) != len(wf.Nodes) {
			t.Errorf("Expected %d nodes, got %d", len(wf.Nodes), len(order))
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
		workflows, err := parser.LoadAllWorkflows()
		if err != nil {
			t.Fatalf("Failed to load all workflows: %v", err)
		}

		if len(workflows) == 0 {
			t.Error("Expected at least one workflow")
		}

		if _, exists := workflows["fix"]; !exists {
			t.Error("Expected fix workflow to exist")
		}
	})

	t.Run("ListWorkflows", func(t *testing.T) {
		infos, err := parser.ListWorkflows()
		if err != nil {
			t.Fatalf("Failed to list workflows: %v", err)
		}

		if len(infos) == 0 {
			t.Error("Expected at least one workflow info")
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
		wf, err := parser.GetWorkflowByType("fix")
		if err != nil {
			t.Fatalf("Failed to get workflow by type: %v", err)
		}

		if wf.WorkflowType != "fix" {
			t.Errorf("Expected workflow_type 'fix', got '%s'", wf.WorkflowType)
		}
	})

	t.Run("DetectCircularDependencies", func(t *testing.T) {
		// Create a workflow with circular dependency
		wf := &types.Workflow{
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

		err := parser.CheckCircularDependencies(wf)
		if err == nil {
			t.Error("Expected error for circular dependency")
		}
	})
}

// TestNodeValidation 测试节点验证
func TestNodeValidation(t *testing.T) {
	parser := workflow.NewParser(getWorkflowsPath(t))

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
