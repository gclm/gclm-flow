# 资源嵌入系统

## 概述

gclm-engine 使用 Go 1.16+ 的 `embed` 包将资源文件嵌入到二进制文件中，实现零依赖部署。

```
二进制文件 = 代码 + migrations + workflows + config
```

---

## 嵌入资源

### embed.go (gclm-engine/embed.go)

```go
//go:embed migrations/*.sql workflows/*.yaml gclm_engine_config.yaml
var AssetsFS embed.FS
```

**当前嵌入**:
- `migrations/*.sql` - 数据库迁移文件 (Goose)
- `workflows/*.yaml` - 默认工作流定义
- `gclm_engine_config.yaml` - 默认配置文件

**未来扩展** (预留):
- `web/static/*` - Web UI 静态文件
- `web/templates/*` - HTML 模板

---

## 资源管理 (internal/assets/embed.go)

### 主要函数

```go
// Init 初始化 assets 包
func Init(fs embed.FS)

// GetFS 返回嵌入的文件系统
func GetFS() *embed.FS

// MigrationsFS 返回 migrations 子文件系统
func MigrationsFS() fs.FS

// WorkflowsFS 返回 workflows 子文件系统
func WorkflowsFS() fs.FS

// GetDefaultConfig 返回默认配置内容
func GetDefaultConfig() ([]byte, error)

// ExportDefaultConfig 导出默认配置到目标目录
func ExportDefaultConfig(targetDir string, force bool) (bool, error)

// ExportBuiltinWorkflows 导出内置工作流到目标目录
func ExportBuiltinWorkflows(targetDir string, force bool) ([]string, error)
```

---

## 数据库迁移 (Goose)

### 迁移文件格式

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS workflows (
    name TEXT PRIMARY KEY,
    ...
);

-- +goose Down
DROP TABLE IF EXISTS workflows;
```

### 迁移历史

| 版本 | 文件 | 说明 |
|:---|:---|:---|
| 00001 | `baseline.sql` | 初始数据库结构 |
| 00002 | `rename_pipeline_to_workflow.sql` | 重命名 pipeline → workflow |
| 00003 | `optimize_indexes.sql` | 优化索引和添加视图 |

### 迁移执行流程

```
1. 设置嵌入的 migrations 文件系统
   ↓
2. 配置 Goose 方言 (sqlite3)
   ↓
3. 执行 goose.Up() 应用未执行的迁移
   ↓
4. 记录迁移版本到 goose_schema_version 表
```

---

## 草稿/正式分离模型

### 工作流生命周期

```
┌─────────────────┐
│  内置 (嵌入)     │ gclm-engine/workflows/*.yaml
│  Builtin        │
└────────┬────────┘
         │ init/export
         ↓
┌─────────────────┐
│  草稿 (可编辑)   │ ~/.gclm-flow/workflows/*.yaml
│  Draft          │
└────────┬────────┘
         │ sync
         ↓
┌─────────────────┐
│  正式 (数据库)   │ workflows 表
│  Production     │
└─────────────────┘
```

### 同步命令

```bash
# 初始化 (内置 → 草稿)
gclm-engine init

# 同步 (草稿 → 正式)
gclm-engine workflow sync                           # 同步所有
gclm-engine workflow sync workflows/feat.yaml      # 同步单个
```

### WorkflowRepository 操作

```go
// InitializeBuiltinWorkflows 从草稿目录加载到数据库
func (r *WorkflowRepository) InitializeBuiltinWorkflows(workflowsDir string) error

// GetWorkflow 从数据库获取
func (r *WorkflowRepository) GetWorkflow(name string) (*WorkflowRecord, error)

// ListWorkflows 列出所有工作流
func (r *WorkflowRepository) ListWorkflows() ([]WorkflowRecord, error)
```

---

## 配置导出

### 自动初始化 (autoInitialize)

```go
func autoInitialize(configDir string) error {
    // 1. 导出默认配置
    assets.ExportDefaultConfig(configDir, false)

    // 2. 导出内置工作流
    workflowsDir := filepath.Join(configDir, "workflows")
    assets.ExportBuiltinWorkflows(workflowsDir, false)

    // 3. 数据库初始化在 db.New() 中完成

    return nil
}
```

### 首次运行检测

```go
func checkNeedsInit(configDir string) bool {
    configFile := filepath.Join(configDir, "gclm_engine_config.yaml")
    workflowsDir := filepath.Join(configDir, "workflows")

    // 检查配置文件和工作流目录是否存在
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        return true
    }

    entries, err := os.ReadDir(workflowsDir)
    if err != nil || len(entries) == 0 {
        return true
    }

    return false
}
```

---

## 零依赖部署优势

### 传统方式 vs 嵌入方式

| 方式 | 优点 | 缺点 |
|:---|:---|:---|
| **传统** | 灵活修改 | 需要分发多个文件 |
| **嵌入** | 零依赖部署 | 修改需重新编译 |

### 嵌入方式优势

1. **单文件部署**: 只需一个二进制文件
2. **版本一致**: 资源与代码版本同步
3. **简化安装**: 无需额外配置文件
4. **跨平台**: 构建时自动嵌入对应平台资源

### 草稿/正式分离优势

1. **可编辑性**: 草稿目录可自由编辑
2. **版本控制**: 草稿文件可纳入 git
3. **可控发布**: sync 命令控制何时发布
4. **回滚能力**: 数据库保留多个版本

---

## 文件位置

| 类型 | 位置 | 说明 |
|:---|:---|:---|
| **嵌入资源** | `gclm-engine/embed.go` | 编译时嵌入 |
| **资源管理** | `internal/assets/embed.go` | 运行时访问 |
| **草稿目录** | `~/.gclm-flow/workflows/` | 用户可编辑 |
| **正式存储** | `workflows` 表 | 数据库 |
| **迁移文件** | `gclm-engine/migrations/` | Goose 格式 |
