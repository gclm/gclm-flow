package types

// NodeConfig represents additional configuration for a node
type NodeConfig map[string]interface{}

// WorkflowNode represents a single node in the workflow
type WorkflowNode struct {
	Ref          string     `yaml:"ref" json:"ref"`
	DisplayName  string     `yaml:"display_name" json:"displayName"`
	Agent        string     `yaml:"agent" json:"agent"`
	Model        string     `yaml:"model" json:"model"`
	Timeout      int        `yaml:"timeout" json:"timeout"` // in seconds
	Required     bool       `yaml:"required" json:"required"`
	DependsOn    []string   `yaml:"depends_on,omitempty" json:"dependsOn,omitempty"`
	ParallelGroup string   `yaml:"parallel_group,omitempty" json:"parallelGroup,omitempty"`
	Config       NodeConfig `yaml:"config,omitempty" json:"config,omitempty"`
}

// CompletionConfig represents workflow completion configuration
type CompletionConfig struct {
	Signal      string `yaml:"signal" json:"signal"`
	FinalStatus string `yaml:"final_status" json:"finalStatus"`
}

// ErrorHandlingConfig represents error handling configuration
type ErrorHandlingConfig struct {
	MaxRetries            int      `yaml:"max_retries" json:"maxRetries"`
	RetryOn               []string `yaml:"retry_on,omitempty" json:"retryOn,omitempty"`
	ContinueOnNonRequired bool     `yaml:"continue_on_non_required" json:"continueOnNonRequired"`
}

// Workflow represents a complete workflow definition
type Workflow struct {
	Name          string              `yaml:"name" json:"name"`
	DisplayName   string              `yaml:"display_name" json:"displayName"`
	Description   string              `yaml:"description,omitempty" json:"description,omitempty"`
	Version       string              `yaml:"version" json:"version"`
	Author        string              `yaml:"author,omitempty" json:"author,omitempty"`
	WorkflowType  string              `yaml:"workflow_type" json:"workflowType"`
	Nodes         []WorkflowNode      `yaml:"nodes" json:"nodes"`
	Completion    CompletionConfig     `yaml:"completion,omitempty" json:"completion,omitempty"`
	ErrorHandling ErrorHandlingConfig `yaml:"error_handling,omitempty" json:"errorHandling,omitempty"`
	// YAML source content (for display)
	ConfigYAML   string              `yaml:"-" json:"configYaml,omitempty"`
	IsBuiltin     bool                `yaml:"-" json:"isBuiltin,omitempty"`
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
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	WorkflowType string `json:"workflowType"`
	ConfigYAML   string `json:"configYaml,omitempty"`
	IsBuiltin     bool   `json:"isBuiltin,omitempty"`
}
