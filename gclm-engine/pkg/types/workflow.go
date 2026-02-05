package types

// NodeConfig represents additional configuration for a node
type NodeConfig map[string]interface{}

// WorkflowNode represents a single node in the workflow
type WorkflowNode struct {
	Ref          string      `yaml:"ref" json:"ref"`
	DisplayName  string      `yaml:"display_name" json:"display_name"`
	Agent        string      `yaml:"agent" json:"agent"`
	Model        string      `yaml:"model" json:"model"`
	Timeout      int         `yaml:"timeout" json:"timeout"` // in seconds
	Required     bool        `yaml:"required" json:"required"`
	DependsOn    []string    `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	ParallelGroup string     `yaml:"parallel_group,omitempty" json:"parallel_group,omitempty"`
	Config       NodeConfig  `yaml:"config,omitempty" json:"config,omitempty"`
}

// CompletionConfig represents workflow completion configuration
type CompletionConfig struct {
	Signal      string `yaml:"signal" json:"signal"`
	FinalStatus string `yaml:"final_status" json:"final_status"`
}

// ErrorHandlingConfig represents error handling configuration
type ErrorHandlingConfig struct {
	MaxRetries            int      `yaml:"max_retries" json:"max_retries"`
	RetryOn               []string `yaml:"retry_on,omitempty" json:"retry_on,omitempty"`
	ContinueOnNonRequired bool     `yaml:"continue_on_non_required" json:"continue_on_non_required"`
}

// Workflow represents a complete workflow definition
type Workflow struct {
	Name          string              `yaml:"name" json:"name"`
	DisplayName   string              `yaml:"display_name" json:"display_name"`
	Description   string              `yaml:"description,omitempty" json:"description,omitempty"`
	Version       string              `yaml:"version" json:"version"`
	Author        string              `yaml:"author,omitempty" json:"author,omitempty"`
	WorkflowType  string              `yaml:"workflow_type" json:"workflow_type"`
	Nodes         []WorkflowNode      `yaml:"nodes" json:"nodes"`
	Completion    CompletionConfig     `yaml:"completion,omitempty" json:"completion,omitempty"`
	ErrorHandling ErrorHandlingConfig `yaml:"error_handling,omitempty" json:"error_handling,omitempty"`
}

// NodeExecutionOrder represents a node with its execution order
type NodeExecutionOrder struct {
	Node     *WorkflowNode
	Order    int
	Parallel int // >0 indicates parallel group number
}

// WorkflowInfo represents basic workflow information
type WorkflowInfo struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	WorkflowType string `json:"workflow_type"`
}
