package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// workflowLoader 实现 domain.WorkflowLoader 接口
// 提供带缓存的工作流加载功能
type workflowLoader struct {
	parser *workflow.Parser
	cache  *loaderCache
	ttl    time.Duration
}

// loaderCache 工作流缓存
type loaderCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
	ttl   time.Duration
}

type cacheItem struct {
	workflow  *types.Workflow
	timestamp time.Time
}

// NewWorkflowLoader 创建新的工作流加载器
func NewWorkflowLoader(parser *workflow.Parser) domain.WorkflowLoader {
	return &workflowLoader{
		parser: parser,
		cache: &loaderCache{
			items: make(map[string]*cacheItem),
			ttl:   5 * time.Minute,
		},
		ttl: 5 * time.Minute,
	}
}

// Load 加载工作流配置（带缓存）
func (l *workflowLoader) Load(ctx context.Context, name string) (*types.Workflow, error) {
	// 检查缓存
	if item := l.cache.get(name); item != nil {
		return item, nil
	}

	// 从文件加载
	workflow, err := l.parser.LoadWorkflow(name)
	if err != nil {
		return nil, fmt.Errorf("load workflow '%s': %w", name, err)
	}

	// 存入缓存
	l.cache.set(name, workflow)

	return workflow, nil
}

// LoadAll 加载所有可用工作流
func (l *workflowLoader) LoadAll(ctx context.Context) (map[string]*types.Workflow, error) {
	return l.parser.LoadAllWorkflows()
}

// Validate 验证工作流配置
func (l *workflowLoader) Validate(ctx context.Context, workflow *types.Workflow) error {
	return l.parser.ValidateWorkflow(workflow)
}

// GetExecutionOrder 计算工作流的执行顺序
func (l *workflowLoader) GetExecutionOrder(ctx context.Context, workflow *types.Workflow) ([]*types.NodeExecutionOrder, error) {
	return l.parser.CalculateExecutionOrder(workflow)
}

// ============================================================================
// loaderCache 方法
// ============================================================================

func (c *loaderCache) get(name string) *types.Workflow {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[name]
	if !exists {
		return nil
	}

	// 检查是否过期
	if time.Since(item.timestamp) > c.ttl {
		delete(c.items, name)
		return nil
	}

	return item.workflow
}

func (c *loaderCache) set(name string, workflow *types.Workflow) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[name] = &cacheItem{
		workflow:  workflow,
		timestamp: time.Now(),
	}
}

func (c *loaderCache) invalidate(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, name)
}

func (c *loaderCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
}
