package pipeline

import (
	"fmt"
	"os"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"gopkg.in/yaml.v3"
)

// Parser handles parsing and validation of pipeline configurations
type Parser struct {
	configDir string
}

// NewParser creates a new pipeline parser
func NewParser(configDir string) *Parser {
	return &Parser{
		configDir: configDir,
	}
}

// LoadPipeline loads a pipeline configuration from a YAML file
func (p *Parser) LoadPipeline(name string) (*types.Pipeline, error) {
	path := fmt.Sprintf("%s/%s.yaml", p.configDir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pipeline file: %w", err)
	}

	var pipeline types.Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to parse pipeline YAML: %w", err)
	}

	// Validate pipeline
	if err := p.ValidatePipeline(&pipeline); err != nil {
		return nil, fmt.Errorf("pipeline validation failed: %w", err)
	}

	return &pipeline, nil
}

// LoadAllPipelines loads all pipeline configurations from the config directory
func (p *Parser) LoadAllPipelines() (map[string]*types.Pipeline, error) {
	entries, err := os.ReadDir(p.configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	pipelines := make(map[string]*types.Pipeline)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		pipeline, err := p.LoadPipeline(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load pipeline %s: %w", name, err)
		}

		pipelines[name] = pipeline
	}

	return pipelines, nil
}

// ValidatePipeline validates a pipeline configuration
func (p *Parser) ValidatePipeline(pipeline *types.Pipeline) error {
	if pipeline.Name == "" {
		return fmt.Errorf("pipeline name is required")
	}

	if pipeline.DisplayName == "" {
		return fmt.Errorf("pipeline display_name is required")
	}

	if pipeline.Version == "" {
		return fmt.Errorf("pipeline version is required")
	}

	if pipeline.WorkflowType == "" {
		return fmt.Errorf("pipeline workflow_type is required")
	}

	if len(pipeline.Nodes) == 0 {
		return fmt.Errorf("pipeline must have at least one node")
	}

	// Validate nodes
	nodeRefs := make(map[string]bool)
	for i, node := range pipeline.Nodes {
		if err := p.ValidateNode(&node, i); err != nil {
			return fmt.Errorf("node %d validation failed: %w", i, err)
		}
		if nodeRefs[node.Ref] {
			return fmt.Errorf("duplicate node ref: %s", node.Ref)
		}
		nodeRefs[node.Ref] = true
	}

	// Validate dependencies
	for _, node := range pipeline.Nodes {
		for _, dep := range node.DependsOn {
			if !nodeRefs[dep] {
				return fmt.Errorf("node %s depends on non-existent node %s", node.Ref, dep)
			}
		}
	}

	// Check for circular dependencies
	if err := p.CheckCircularDependencies(pipeline); err != nil {
		return err
	}

	return nil
}

// ValidateNode validates a single node configuration
func (p *Parser) ValidateNode(node *types.PipelineNode, index int) error {
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

// CheckCircularDependencies checks for circular dependencies in the pipeline
func (p *Parser) CheckCircularDependencies(pipeline *types.Pipeline) error {
	// Build adjacency list
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	for _, node := range pipeline.Nodes {
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

	if processedCount != len(pipeline.Nodes) {
		return fmt.Errorf("circular dependency detected in pipeline")
	}

	return nil
}

// CalculateExecutionOrder calculates the execution order for pipeline nodes
// Returns nodes in topological order with parallel group information
func (p *Parser) CalculateExecutionOrder(pipeline *types.Pipeline) ([]*types.NodeExecutionOrder, error) {
	// Build dependency map and in-degree count
	nodeMap := make(map[string]*types.PipelineNode)
	inDegree := make(map[string]int)
	dependents := make(map[string][]string) // reverse graph

	for i := range pipeline.Nodes {
		node := &pipeline.Nodes[i]
		nodeMap[node.Ref] = node
		inDegree[node.Ref] = len(node.DependsOn)
	}

	// Build reverse graph
	for _, node := range pipeline.Nodes {
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
	if len(result) != len(pipeline.Nodes) {
		return nil, fmt.Errorf("cycle detected in pipeline dependencies")
	}

	return result, nil
}

// GetPipelineByWorkflowType finds a pipeline by workflow type
func (p *Parser) GetPipelineByWorkflowType(workflowType string) (*types.Pipeline, error) {
	pipelines, err := p.LoadAllPipelines()
	if err != nil {
		return nil, err
	}

	for _, pipeline := range pipelines {
		if pipeline.WorkflowType == string(workflowType) {
			return pipeline, nil
		}
	}

	return nil, fmt.Errorf("no pipeline found for workflow type: %s", workflowType)
}

// ListPipelines returns a list of all available pipeline information
func (p *Parser) ListPipelines() ([]*types.PipelineInfo, error) {
	pipelines, err := p.LoadAllPipelines()
	if err != nil {
		return nil, err
	}

	info := make([]*types.PipelineInfo, 0, len(pipelines))
	for _, p := range pipelines {
		info = append(info, &types.PipelineInfo{
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Version:     p.Version,
			WorkflowType: p.WorkflowType,
		})
	}

	return info, nil
}

// LoadYAMLFile loads a pipeline from any YAML file path
func (p *Parser) LoadYAMLFile(yamlFile string) (*types.Pipeline, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var pipeline types.Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &pipeline, nil
}

// PipelineToYAML converts a pipeline to YAML format
func (p *Parser) PipelineToYAML(pipeline *types.Pipeline) ([]byte, error) {
	data, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to YAML: %w", err)
	}

	return data, nil
}
