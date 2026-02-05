# 目录结构

## 项目结构

```
gclm-flow/
├── README.md                  # 项目文档
├── CLAUDE.md                 # 项目特定配置
├── CLAUDE.example.md         # 全局配置模板 (→ ~/.claude/CLAUDE.md)
├── install.sh                # 安装脚本 (gclm-engine + 全局配置)
├── settings.example.json      # MCP 配置示例
│
├── gclm-engine/              # Go 引擎 (工作流编排和状态管理)
│   ├── main.go               # 入口文件
│   ├── go.mod                # Go 模块定义
│   ├── Makefile              # 构建脚本
│   ├── internal/             # 内部包
│   │   ├── cli/              # CLI 命令 (cobra)
│   │   │   └── commands.go   # workflow/task 命令
│   │   ├── db/               # 数据库操作
│   │   │   ├── database.go   # Database 初始化
│   │   │   └── workflow.go   # Workflow Repository
│   │   ├── pipeline/         # YAML 解析
│   │   │   └── parser.go     # 工作流解析、依赖检查
│   │   └── service/          # 业务逻辑
│   │       ├── task.go       # 任务服务 (智能分流)
│   │       └── errors.go     # 错误定义
│   ├── pkg/                  # 共享包
│   │   └── types/            # 类型定义
│   │       ├── pipeline.go   # Pipeline 类型
│   │       └── types.go      # 通用类型
│   └── test/                 # 测试文件
│
├── workflows/                # 工作流定义 (统一位置)
│   ├── code_simple.yaml     # CODE_SIMPLE 工作流
│   ├── code_complex.yaml    # CODE_COMPLEX 工作流
│   ├── document.yaml        # DOCUMENT 工作流
│   ├── examples/            # 自定义工作流示例
│   │   ├── custom_simple.yaml
│   │   ├── custom_document.yaml
│   │   └── custom_complex.yaml
│   └── README.md            # 工作流说明
│
├── agents/                   # Agent 定义
│   ├── investigator.md      # 代码库调查 Agent
│   ├── architect.md         # 架构设计 Agent
│   ├── spec-guide.md        # SpecDD 指导 Agent
│   ├── tdd-guide.md         # TDD 指导 Agent
│   ├── worker.md            # 任务执行 Agent
│   ├── code-reviewer.md     # 代码审查 Agent
│   └── recorder.md          # 文档记录 Agent
│
├── skills/                   # Claude Code Skills
│   └── gclm/
│       └── SKILL.md         # 智能分流工作流 Skill
│
├── rules/                    # 工作流规则文档
│   ├── agents.md            # Agent 编排规则
│   ├── llmdoc.md            # llmdoc 文档规则
│   ├── phases.md            # 阶段流程规则
│   ├── spec.md              # SpecDD 规范
│   └── tdd.md               # TDD 规范
│
├── commands/                 # Claude Code 命令定义 (已废弃)
│   ├── gclm.md
│   ├── investigate.md
│   ├── tdd.md
│   ├── spec.md
│   └── llmdoc.md
│
├── hooks/                    # Claude Code Hooks
│   └── stop-gclm-loop.sh    # 停止 Hook
│
├── .github/                 # GitHub 配置
│   └── workflows/
│       └── release.yml      # Release 工作流
│
└── llmdoc/                   # LLM 优化文档
    ├── index.md             # 文档索引
    ├── overview/            # 项目概览
    │   ├── project.md
    │   ├── tech-stack.md
    │   └── structure.md
    ├── architecture/        # 架构设计
    │   ├── system.md        # 系统架构
    │   ├── workflows.md     # 工作流配置
    │   ├── agents.md        # Agent 体系
    │   ├── code-search.md   # 代码搜索策略
    │   └── database.md      # 数据库设计
    ├── guides/              # 使用指南
    │   ├── installation.md
    │   ├── quickstart.md
    │   └── workflow-development.md
    └── reference/           # 参考文档
        ├── commands.md
        ├── workflows.md
        └── configuration.md
```

---

## 关键目录说明

### `gclm-engine/`

Go 引擎，负责工作流编排和状态管理。

| 子目录 | 职责 |
|:---|:---|
| `internal/cli/` | CLI 命令，输出 JSON 供 Skills 调用 |
| `internal/db/` | SQLite 数据库操作 |
| `internal/pipeline/` | YAML 工作流解析、依赖检查 |
| `internal/service/` | 任务服务（智能分流、阶段管理） |
| `pkg/types/` | 共享数据结构定义 |

### `workflows/`

工作流定义目录，所有工作流 YAML 文件统一存放。

| 文件 | 类型 | 阶段数 |
|:---|:---|:---:|
| `code_simple.yaml` | CODE_SIMPLE | 6 |
| `code_complex.yaml` | CODE_COMPLEX | 9 |
| `document.yaml` | DOCUMENT | 7 |
| `examples/*.yaml` | 自定义 | 可变 |

### `agents/`

Agent 定义文件，定义每个 Agent 的职责、模型选择和使用场景。

### `skills/gclm/`

主 Skill 定义，编排工作流调用 gclm-engine 命令。

### `rules/`

工作流规则文档，定义 Agent 编排、阶段流程、SpecDD/TDD 规范。

### `.github/workflows/`

GitHub Actions 工作流，用于 CI/CD 和 Release。

### `llmdoc/`

LLM 优化的项目文档，帮助 AI 理解项目结构。

---

## 安装位置

install.sh 安装后的文件位置：

| 文件/目录 | 安装位置 |
|:---|:---|
| gclm-engine 二进制 | `~/.gclm-flow/gclm-engine` |
| 工作流文件 | `~/.gclm-flow/workflows/*.yaml` |
| SQLite 数据库 | `~/.gclm-flow/gclm-engine.db` |
| 全局配置 | `~/.claude/CLAUDE.md` |
| Agents | `~/.claude/agents/*.md` |
| Skills | `~/.claude/skills/gclm/SKILL.md` |
| Rules | `~/.claude/rules/*.md` |
| Hooks | `~/.claude/hooks/*.sh` |
