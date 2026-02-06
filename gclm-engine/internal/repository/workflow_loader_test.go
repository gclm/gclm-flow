package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TestWorkflowLoader 测试 WorkflowLoader 适配器
func TestWorkflowLoader(t *testing.T) {
	t.Run("LoadWorkflow", func(t *testing.T) {
		workflowsDir := getWorkflowsDir(t)

		parser := workflow.NewParser(workflowsDir)
		loader := NewWorkflowLoader(parser)

		ctx := context.Background()

		// Test loading workflow by name
		wf, err := loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		if wf.Name != "code_simple" {
			t.Errorf("Expected name 'code_simple', got '%s'", wf.Name)
		}

		if len(wf.Nodes) == 0 {
			t.Error("Expected at least one node")
		}
	})

	t.Run("LoadAllWorkflows", func(t *testing.T) {
		workflowsDir := getWorkflowsDir(t)

		parser := workflow.NewParser(workflowsDir)
		loader := NewWorkflowLoader(parser)

		ctx := context.Background()

		// Load all workflows
		workflows, err := loader.LoadAll(ctx)
		if err != nil {
			t.Fatalf("Failed to load all workflows: %v", err)
		}

		if len(workflows) == 0 {
			t.Error("Expected at least one workflow")
		}

		// Check that required fields are present
		for name, wf := range workflows {
			if name == "" {
				t.Error("Workflow name is required")
			}
			if wf.Name == "" {
				t.Error("Workflow name field is required")
			}
		}
	})

	t.Run("GetExecutionOrder", func(t *testing.T) {
		workflowsDir := getWorkflowsDir(t)

		parser := workflow.NewParser(workflowsDir)
		loader := NewWorkflowLoader(parser)

		ctx := context.Background()

		// Load workflow
		wf, err := loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		// Get execution order
		order, err := loader.GetExecutionOrder(ctx, wf)
		if err != nil {
			t.Fatalf("Failed to get execution order: %v", err)
		}

		if len(order) == 0 {
			t.Error("Expected at least one node in execution order")
		}

		if len(order) != len(wf.Nodes) {
			t.Errorf("Expected %d nodes, got %d", len(wf.Nodes), len(order))
		}
	})

	t.Run("ValidateWorkflow", func(t *testing.T) {
		workflowsDir := getWorkflowsDir(t)

		parser := workflow.NewParser(workflowsDir)
		loader := NewWorkflowLoader(parser)

		ctx := context.Background()

		// Load workflow
		wf, err := loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		// Validate workflow
		err = loader.Validate(ctx, wf)
		if err != nil {
			t.Errorf("Expected valid workflow, got error: %v", err)
		}
	})

	t.Run("Caching", func(t *testing.T) {
		workflowsDir := getWorkflowsDir(t)

		parser := workflow.NewParser(workflowsDir)
		loader := NewWorkflowLoader(parser)

		ctx := context.Background()

		// First call - cache miss
		start1 := time.Now()
		_, err := loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("First call failed: %v", err)
		}
		duration1 := time.Since(start1)

		// Second call - cache hit
		start2 := time.Now()
		_, err = loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("Second call failed: %v", err)
		}
		duration2 := time.Since(start2)

		// Cache hit should be faster
		// Note: This is a weak test, but demonstrates caching behavior
		// In a real scenario, you'd use a more precise timing method
		t.Logf("First call: %v, Second call: %v", duration1, duration2)
	})
}

// TestWorkflowLoaderInvalidInput 测试无效输入
func TestWorkflowLoaderInvalidInput(t *testing.T) {
	workflowsDir := getWorkflowsDir(t)

	parser := workflow.NewParser(workflowsDir)
	loader := NewWorkflowLoader(parser)

	ctx := context.Background()

	t.Run("LoadNotFound", func(t *testing.T) {
		_, err := loader.Load(ctx, "nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent workflow")
		}
	})

	t.Run("ValidateInvalidWorkflow", func(t *testing.T) {
		// Create an invalid workflow (empty nodes)
		wf := &types.Workflow{
			Name:         "invalid",
			DisplayName:  "Invalid Workflow",
			Version:      "0.1.0",
			WorkflowType: "fix",
			Nodes:        []types.WorkflowNode{}, // Empty nodes
		}

		err := loader.Validate(ctx, wf)
		if err == nil {
			t.Error("Expected error for workflow with no nodes")
		}
	})
}

// TestWorkflowLoaderExecutionOrderValidation 测试执行顺序计算
func TestWorkflowLoaderExecutionOrderValidation(t *testing.T) {
	workflowsDir := getWorkflowsDir(t)

	parser := workflow.NewParser(workflowsDir)
	loader := NewWorkflowLoader(parser)

	ctx := context.Background()

	t.Run("ExecutionOrderSequence", func(t *testing.T) {
		wf, err := loader.Load(ctx, "code_simple")
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}

		order, err := loader.GetExecutionOrder(ctx, wf)
		if err != nil {
			t.Fatalf("Failed to get execution order: %v", err)
		}

		// Check that order field is set correctly
		for i, nodeOrder := range order {
			if nodeOrder.Order != i {
				t.Errorf("Expected order %d, got %d", i, nodeOrder.Order)
			}
		}

		// Check that dependencies are respected
		// Find discovery and clarification
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
}

// getWorkflowsDir finds the workflows directory
func getWorkflowsDir(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Find workflows directory
	workflowsDir := "workflows"
	for i := 0; i < 4; i++ {
		testPath := filepath.Join(wd, workflowsDir)
		if info, err := os.Stat(testPath); err == nil && info.IsDir() {
			// Check if it has YAML files
			entries, err := os.ReadDir(testPath)
			if err == nil && len(entries) > 0 {
				return testPath
			}
		}
		wd = filepath.Dir(wd)
		if filepath.Base(wd) == "gclm-engine" {
			return filepath.Join(wd, "workflows")
		}
	}

	t.Fatal("Could not find workflows directory")
	return ""
}
