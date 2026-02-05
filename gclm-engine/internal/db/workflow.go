package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/config"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"gopkg.in/yaml.v3"
)

// WorkflowRepository handles workflow CRUD operations
type WorkflowRepository struct {
	db *Database
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *Database) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

// WorkflowRecord represents a workflow in the database
type WorkflowRecord struct {
	Name        string
	DisplayName string
	Description string
	WorkflowType string
	Version     string
	IsBuiltin   bool
	ConfigYAML  string
}

// InitializeBuiltinWorkflows loads builtin workflows from YAML files into the database
// Scans the workflows directory for all .yaml files (excluding examples/)
// Note: Directory should already exist (created by autoInitialize or `init` command)
func (r *WorkflowRepository) InitializeBuiltinWorkflows(workflowsDir string) error {
	// Check if directory exists
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist - this is an error since autoInitialize should have created it
			return fmt.Errorf("workflows directory not found: %s (run 'gclm-engine workflow init' to setup)", workflowsDir)
		}
		return fmt.Errorf("failed to read workflows directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories and non-YAML files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		// Skip files in examples/ subdirectory (if processed recursively)
		if strings.HasPrefix(entry.Name(), "example") {
			continue
		}

		yamlPath := filepath.Join(workflowsDir, entry.Name())

		// Read YAML file
		yamlData, err := os.ReadFile(yamlPath)
		if err != nil {
			return fmt.Errorf("failed to read workflow file %s: %w", yamlPath, err)
		}

		// Parse YAML to get workflow metadata
		var wf types.Workflow
		if err := yaml.Unmarshal(yamlData, &wf); err != nil {
			return fmt.Errorf("failed to parse workflow YAML %s: %w", yamlPath, err)
		}

		// Extract workflow name from filename (remove .yaml extension)
		workflowName := strings.TrimSuffix(entry.Name(), ".yaml")

		// Check if workflow already exists
		var exists bool
		checkErr := r.db.conn.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM workflows WHERE name = ? AND is_builtin = 1)",
			workflowName,
		).Scan(&exists)
		if checkErr != nil {
			return fmt.Errorf("failed to check workflow existence: %w", checkErr)
		}

		// Skip if already exists
		if exists {
			continue
		}

		// Get workflow_type from workflow (required)
		workflowType := wf.WorkflowType
		if workflowType == "" {
			return fmt.Errorf("workflow %s: workflow_type is required", workflowName)
		}

		// Validate workflow_type against config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.ValidateWorkflowType(workflowType); err != nil {
			return fmt.Errorf("workflow %s: %w", workflowName, err)
		}

		// Insert into database
		_, err = r.db.conn.Exec(`
			INSERT INTO workflows (name, display_name, description, workflow_type, version, is_builtin, config_yaml)
			VALUES (?, ?, ?, ?, ?, 1, ?)
		`, workflowName, wf.DisplayName, wf.Description, workflowType, wf.Version, string(yamlData))

		if err != nil {
			return fmt.Errorf("failed to insert workflow %s: %w", workflowName, err)
		}
	}

	return nil
}

// ListWorkflows returns all workflows
func (r *WorkflowRepository) ListWorkflows() ([]WorkflowRecord, error) {
	rows, err := r.db.conn.Query(`
		SELECT name, display_name, description, workflow_type, version, is_builtin, config_yaml
		FROM workflows
		ORDER BY is_builtin DESC, name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	defer rows.Close()

	var workflows []WorkflowRecord
	for rows.Next() {
		var w WorkflowRecord
		err := rows.Scan(&w.Name, &w.DisplayName, &w.Description, &w.WorkflowType, &w.Version, &w.IsBuiltin, &w.ConfigYAML)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow: %w", err)
		}
		workflows = append(workflows, w)
	}

	return workflows, nil
}

// GetWorkflow retrieves a workflow by name
func (r *WorkflowRepository) GetWorkflow(name string) (*WorkflowRecord, error) {
	var w WorkflowRecord
	err := r.db.conn.QueryRow(`
		SELECT name, display_name, description, workflow_type, version, is_builtin, config_yaml
		FROM workflows
		WHERE name = ?
	`, name).Scan(&w.Name, &w.DisplayName, &w.Description, &w.WorkflowType, &w.Version, &w.IsBuiltin, &w.ConfigYAML)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return &w, nil
}

// GetWorkflowByType retrieves a workflow by workflow_type
func (r *WorkflowRepository) GetWorkflowByType(workflowType string) (*WorkflowRecord, error) {
	var w WorkflowRecord
	err := r.db.conn.QueryRow(`
		SELECT name, display_name, description, workflow_type, version, is_builtin, config_yaml
		FROM workflows
		WHERE workflow_type = ?
		LIMIT 1
	`, workflowType).Scan(&w.Name, &w.DisplayName, &w.Description, &w.WorkflowType, &w.Version, &w.IsBuiltin, &w.ConfigYAML)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow of type '%s' not found", workflowType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow by type: %w", err)
	}

	return &w, nil
}

// InstallWorkflow installs a new workflow from YAML content
func (r *WorkflowRepository) InstallWorkflow(name string, yamlData []byte) error {
	// Parse YAML to get metadata
	var wf types.Workflow
	if err := yaml.Unmarshal(yamlData, &wf); err != nil {
		return fmt.Errorf("failed to parse workflow YAML: %w", err)
	}

	// Check if workflow already exists
	var exists bool
	err := r.db.conn.QueryRow("SELECT EXISTS(SELECT 1 FROM workflows WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check workflow existence: %w", err)
	}

	if exists {
		return fmt.Errorf("workflow '%s' already exists", name)
	}

	// Determine workflow type from workflow
	workflowType := wf.WorkflowType
	if workflowType == "" {
		// Fallback: infer from name
		switch name {
		case "document":
			workflowType = "document"
		case "code_simple", "bug_fix_fast":
			workflowType = "code_simple"
		default:
			workflowType = "code_complex"
		}
	}

	// Insert into database
	_, err = r.db.conn.Exec(`
		INSERT INTO workflows (name, display_name, description, workflow_type, version, is_builtin, config_yaml)
		VALUES (?, ?, ?, ?, ?, 0, ?)
	`, name, wf.DisplayName, wf.Description, workflowType, wf.Version, string(yamlData))

	if err != nil {
		return fmt.Errorf("failed to install workflow: %w", err)
	}

	return nil
}

// UninstallWorkflow removes a custom workflow
func (r *WorkflowRepository) UninstallWorkflow(name string) error {
	// Check if it's a builtin workflow
	var isBuiltin bool
	err := r.db.conn.QueryRow("SELECT is_builtin FROM workflows WHERE name = ?", name).Scan(&isBuiltin)
	if err == sql.ErrNoRows {
		return fmt.Errorf("workflow '%s' not found", name)
	}
	if err != nil {
		return fmt.Errorf("failed to check workflow: %w", err)
	}

	if isBuiltin {
		return fmt.Errorf("cannot uninstall builtin workflow '%s'", name)
	}

	// Delete the workflow
	result, err := r.db.conn.Exec("DELETE FROM workflows WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("failed to uninstall workflow: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("workflow '%s' not found", name)
	}

	return nil
}
