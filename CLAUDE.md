# gclm-flow 项目配置

本文件为 Claude Code (claude.ai/code) 在此代码库中工作提供指导。

---

## 项目概述

**gclm-flow** 是一个基于 Go 引擎的智能工作流系统，支持自定义工作流 YAML 配置和多 Agent 并行执行。

```
用户请求 → gclm-engine (Go 引擎) → 工作流编排 → Agent 执行
    ↓              ↓                    ↓
 自然语言    SQLite 状态管理      多 Agent 并行
```

---

## 快速开始

### 安装

```bash
# 运行安装脚本（下载二进制 + 同步工作流）
./install.sh

# 或手动安装
cd gclm-engine
make build
make install
```

### 使用工作流

```bash
# 列出所有工作流
gclm-engine workflow list

# 开始工作流（自动检测类型）
gclm-engine workflow start "修复登录页面 bug"

# 查看当前阶段
gclm-engine task current <task-id>

# 完成阶段
gclm-engine task complete <task-id> <phase-id> --output "结果"
```

### 开发

```bash
# 本地开发构建
cd gclm-engine && make dev

# 运行测试
make test
```

---

## 工作流类型

| 类型 | 检测关键词 | 适用场景 |
|:---|:---|:---|
| 📝 **DOCUMENT** | 文档、方案、设计、需求 | 文档编写 |
| 🔧 **CODE_SIMPLE** | bug、修复、error | Bug修复/小修改 |
| 🚀 **CODE_COMPLEX** | 功能、模块、开发、重构 | 新功能/复杂变更 |

---

## 核心组件

| 组件 | 位置 | 用途 |
|:---|:---|:---|
| **Go 引擎 CLI** | `gclm-engine/internal/cli/` | 命令接口，为 skills 提供 JSON 输出 |
| **任务服务** | `gclm-engine/internal/service/` | 核心工作流逻辑，阶段转换 |
| **数据库层** | `gclm-engine/internal/db/` | 任务/阶段/事件的 SQLite 持久化 |
| **流水线解析器** | `gclm-engine/internal/pipeline/` | YAML 工作流解析，依赖解析 |
| **工作流 YAML** | `workflows/` | 定义 DOCUMENT、CODE_SIMPLE、CODE_COMPLEX 流程 |
| **Skills** | `skills/gclm/SKILL.md` | 编排工作流的主 skill |
| **Agents** | `agents/*.md` | Agent 定义 (investigator、architect、tdd-guide 等) |

---

## 目录结构

```
gclm-flow/
├── gclm-engine/          # Go 引擎
│   ├── main.go           # 入口文件
│   ├── internal/
│   │   ├── cli/          # CLI 命令 (cobra)
│   │   ├── db/           # SQLite 操作
│   │   ├── pipeline/     # YAML 解析器
│   │   └── service/      # 任务服务 (工作流逻辑)
│   ├── pkg/types/        # 共享类型
│   └── Makefile
├── workflows/            # 工作流定义（统一位置）
│   ├── *.yaml           # 内置工作流
│   └── examples/        # 自定义工作流示例
├── agents/              # Agent 定义
├── skills/              # Skill 定义
├── rules/               # 工作流规则 (phases, tdd, spec)
└── install.sh           # 安装脚本
```

---

## 工作流配置

工作流在 `workflows/` 中通过 YAML 定义：

```yaml
name: code_simple
workflow_type: "CODE_SIMPLE"
nodes:
  - ref: discovery
    display_name: "Discovery / 需求发现"
    agent: investigator
    model: haiku
    timeout: 60
    required: true
  - ref: clarification
    depends_on: [discovery]
    # ... 更多节点
```

### 添加新工作流

1. 在 `workflows/` 中创建 YAML 文件
2. 使用 `depends_on` 定义节点依赖
3. 使用 `parallel_group` 实现并行执行
4. 用 `required: true` 标记关键节点

---

## Skills 集成

主 skill: `skills/gclm/SKILL.md`

**关键集成点：**
- `workflow start <prompt>` → 创建任务，返回第一阶段
- `task current <task-id>` → 获取下一个待执行阶段
- `task complete <task-id> <phase-id> --output "..."` → 标记阶段完成

---

## 约定规范

### 工作流类型检测

关键词评分系统 (位于 `service/task.go`)：
- 文档短语 (+5): "编写文档", "方案设计"
- 文档单词 (+3): "文档", "方案", "需求"
- Bug 修复短语 (-5): "修复bug", "fix bug"
- Bug 修复单词 (-3): "bug", "修复", "debug"
- 功能开发单词 (-1): "功能", "模块", "开发"

阈值：score >= 3 → DOCUMENT, score <= -3 → CODE_SIMPLE, 其他 → CODE_COMPLEX

### 阶段状态

`pending` → `running` → `completed` / `failed` / `skipped`

---

## 数据库结构

位于 `~/.gclm-flow/gclm-engine.db`：

- **tasks**: id, pipeline_id, prompt, workflow_type, status, current_phase, total_phases
- **task_phases**: id, task_id, phase_name, agent, model, status, output_text
- **events**: id, task_id, phase_id, event_type, data (审计日志)

---

## 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
cd gclm-engine && go test ./internal/cli -v
cd gclm-engine && go test ./internal/service -v

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 发布流程

1. 更新 `gclm-engine/internal/cli/commands.go` 中的版本
2. 创建 git 标签: `git tag v0.x.x`
3. 推送标签: `git push origin v0.x.x`
4. GitHub Actions 构建多平台二进制文件

---

## 重要约束

1. **SQLite 单写入者**: 数据库使用 `SetMaxOpenConns(1)`
2. **WAL 模式**: 启用以提升并发性
3. **工作流状态**: 存储在 `~/.gclm-flow/gclm-engine.db`
4. **JSON 输出**: 所有引擎命令支持 `--json` 标志
5. **阶段依赖**: 必须形成 DAG

---

## 依赖项

- `github.com/spf13/cobra` - CLI 框架
- `github.com/mattn/go-sqlite3` - SQLite 驱动 (需要 CGO)
- `gopkg.in/yaml.v3` - YAML 解析
