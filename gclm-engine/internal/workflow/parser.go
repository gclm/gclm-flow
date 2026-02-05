package workflow

import (
	"fmt"
	"os"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"gopkg.in/yaml.v3"
)

// Parser handles parsing and validation of workflow configurations
type Parser struct {
	configDir string
}

// NewParser creates a new workflow parser
func NewParser(configDir string) *Parser {
	return &Parser{
		configDir: configDir,
	}
}

// LoadWorkflow loads a workflow configuration from a YAML file
func (p *Parser) LoadWorkflow(name string) (*types.Workflow, error) {
	path := fmt.Sprintf("%s/%s.yaml", p.configDir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow types.Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow YAML: %w", err)
	}

	// Validate workflow
	if err := p.ValidateWorkflow(&workflow); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	return &workflow, nil
}

// LoadAllWorkflows loads all workflow configurations from the config directory
func (p *Parser) LoadAllWorkflows() (map[string]*types.Workflow, error) {
	entries, err := os.ReadDir(p.configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	workflows := make(map[string]*types.Workflow)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		workflow, err := p.LoadWorkflow(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load workflow %s: %w", name, err)
		}

		workflows[name] = workflow
	}

	return workflows, nil
}

// ValidateWorkflow validates a workflow configuration
func (p *Parser) ValidateWorkflow(workflow *types.Workflow) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if workflow.DisplayName == "" {
		return fmt.Errorf("workflow display_name is required")
	}

	if workflow.Version == "" {
		return fmt.Errorf("workflow version is required")
	}

	if workflow.WorkflowType == "" {
		return fmt.Errorf("workflow workflow_type is required")
	}

	if len(workflow.Nodes) == 0 {
		return fmt.Errorf("workflow must have at least one node")
	}

	// Validate nodes
	nodeRefs := make(map[string]bool)
	for i, node := range workflow.Nodes {
		if err := p.ValidateNode(&node, i); err != nil {
			return fmt.Errorf("node %d validation failed: %w", i, err)
		}
		if nodeRefs[node.Ref] {
			return fmt.Errorf("duplicate node ref: %s", node.Ref)
		}
		nodeRefs[node.Ref] = true
	}

	// Validate dependencies
	for _, node := range workflow.Nodes {
		for _, dep := range node.DependsOn {
			if !nodeRefs[dep] {
				return fmt.Errorf("node %s depends on non-existent node %s", node.Ref, dep)
			}
		}
	}

	// Check for circular dependencies
	if err := p.CheckCircularDependencies(workflow); err != nil {
		return err
	}

	return nil
}

// ValidateNode validates a single node configuration
func (p *Parser) ValidateNode(node *types.WorkflowNode, index int) error {
	if node.Ref == "" {
		return fmt.Errorf("node ref is required")
	}

	if node.DisplayName == "" {
		return fmt.Errorf("node display_name is required")
	}

	if node.Agent == "" {
		return fmt.Errorf("node agent is required")
	}

	if node.Model == "" {
		return fmt.Errorf("node model is required")
	}

	if node.Timeout <= 0 {
		return fmt.Errorf("node timeout must be positive")
	}

	return nil
}

// CheckCircularDependencies checks for circular dependencies in the workflow
func (p *Parser) CheckCircularDependencies(workflow *types.Workflow) error {
	// Build adjacency list
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	for _, node := range workflow.Nodes {
		graph[node.Ref] = node.DependsOn
		inDegree[node.Ref] = len(node.DependsOn)
	}

	// Kahn's algorithm for topological sort
	queue := make([]string, 0)
	visited := make(map[string]bool)

	// Find all nodes with no dependencies
	for ref, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, ref)
		}
	}

	processedCount := 0
	for len(queue) > 0 {
		// Dequeue a node
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		processedCount++

		// Reduce in-degree for all dependent nodes
		for ref, deps := range graph {
			for _, dep := range deps {
				if dep == current {
					inDegree[ref]--
					if inDegree[ref] == 0 {
						queue = append(queue, ref)
					}
				}
			}
		}
	}

	if processedCount != len(workflow.Nodes) {
		return fmt.Errorf("circular dependency detected in workflow")
	}

	return nil
}

// CalculateExecutionOrder calculates the execution order for workflow nodes
// Returns nodes in topological order with parallel group information
func (p *Parser) CalculateExecutionOrder(workflow *types.Workflow) ([]*types.NodeExecutionOrder, error) {
	// Build dependency map and in-degree count
	nodeMap := make(map[string]*types.WorkflowNode)
	inDegree := make(map[string]int)
	dependents := make(map[string][]string) // reverse graph

	for i := range workflow.Nodes {
		node := &workflow.Nodes[i]
		nodeMap[node.Ref] = node
		inDegree[node.Ref] = len(node.DependsOn)
	}

	// Build reverse graph
	for _, node := range workflow.Nodes {
		for _, dep := range node.DependsOn {
			dependents[dep] = append(dependents[dep], node.Ref)
		}
	}

	// Find nodes with no dependencies
	queue := make([]string, 0)
	for ref, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, ref)
		}
	}

	// Process nodes in topological order
	order := 0
	parallelGroup := 0
	result := make([]*types.NodeExecutionOrder, 0)

	for len(queue) > 0 {
		batchSize := len(queue)
		parallelGroup++

		// Process all nodes in current level (can be executed in parallel)
		for i := 0; i < batchSize; i++ {
			ref := queue[0]
			queue = queue[1:]

			node := nodeMap[ref]
			result = append(result, &types.NodeExecutionOrder{
				Node:     node,
				Order:    order,
				Parallel: parallelGroup,
			})
			order++

			// Reduce in-degree for dependent nodes
			for _, depRef := range dependents[ref] {
				inDegree[depRef]--
				if inDegree[depRef] == 0 {
					queue = append(queue, depRef)
				}
			}
		}
	}

	// Verify all nodes were processed
	if len(result) != len(workflow.Nodes) {
		return nil, fmt.Errorf("cycle detected in workflow dependencies")
	}

	return result, nil
}

// GetWorkflowByType finds a workflow by workflow type
func (p *Parser) GetWorkflowByType(workflowType string) (*types.Workflow, error) {
	workflows, err := p.LoadAllWorkflows()
	if err != nil {
		return nil, err
	}

	for _, workflow := range workflows {
		if workflow.WorkflowType == string(workflowType) {
			return workflow, nil
		}
	}

	return nil, fmt.Errorf("no workflow found for workflow type: %s", workflowType)
}

// ListWorkflows returns a list of all available workflow information
func (p *Parser) ListWorkflows() ([]*types.WorkflowInfo, error) {
	workflows, err := p.LoadAllWorkflows()
	if err != nil {
		return nil, err
	}

	info := make([]*types.WorkflowInfo, 0, len(workflows))
	for _, w := range workflows {
		info = append(info, &types.WorkflowInfo{
			Name:         w.Name,
			DisplayName:  w.DisplayName,
			Description:  w.Description,
			Version:      w.Version,
			WorkflowType: w.WorkflowType,
		})
	}

	return info, nil
}

// LoadYAMLFile loads a workflow from any YAML file path
func (p *Parser) LoadYAMLFile(yamlFile string) (*types.Workflow, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var workflow types.Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &workflow, nil
}

// WorkflowToYAML converts a workflow to YAML format
func (p *Parser) WorkflowToYAML(workflow *types.Workflow) ([]byte, error) {
	data, err := yaml.Marshal(workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %w", err)
	}

	return data, nil
}
