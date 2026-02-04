# gclm-flow

> 智能分流工作流插件 - 基于 Go 引擎，自动判断任务类型，选择最优开发流程

---

## 快速导航

| 我想... | 立即跳转 |
|:---|:---|
| **快速安装** | [安装指南](#安装) |
| **怎么用** | [使用方法](#使用方法) |
| **工作流是什么** | [工作流概览](#工作流概览) |
| **有哪些命令** | [可用命令](#可用命令) |
| **项目结构** | [目录结构](#目录结构) |

---

## 安装

```bash
# 克隆仓库
git clone https://github.com/gclm/gclm-flow.git
cd gclm-flow

# 运行安装脚本
bash install.sh
```

**安装脚本会自动**：
1. 编译并安装 `gclm-engine` 到 `~/.gclm-flow/`
2. 安装 agents、skills、commands 到 `~/.claude/`
3. 同步工作流 YAML 文件
4. 初始化数据库

---

## 使用方法

### 基本使用

```bash
# 启动智能工作流 (自动判断任务类型)
/gclm 实现用户登录功能

# 快速代码库调查
/investigate 认证系统是怎么工作的？

# TDD 开发
/tdd 实现密码验证函数

# SpecDD 开发
/spec 生成支付模块规范文档

# 智能提交
/commit
```

### gclm-engine 命令

```bash
# 查看所有工作流
gclm-engine workflow list

# 一键开始工作流
gclm-engine workflow start "修复登录页面 bug"

# 查看任务状态
gclm-engine task current <task-id>

# 列出所有任务
gclm-engine task list
```

---

## 工作流概览

### 三种工作流类型

| 类型 | 触发关键词 | 适用场景 |
|:---|:---|:---|
| 📝 **DOCUMENT** | 文档、方案、设计、需求 | 文档编写、方案设计、架构设计 |
| 🔧 **CODE_SIMPLE** | bug、修复、error、fix | Bug 修复、小修改、单文件变更 |
| 🚀 **CODE_COMPLEX** | 功能、模块、开发、重构 | 新功能、模块开发、跨文件变更 |

### 工作流执行流程

```
用户请求 → gclm-engine 自动检测工作流类型 → 执行阶段 → 完成任务
```

**工作流由 YAML 配置定义**，位于 `workflows/` 目录，支持自定义扩展。

---

## 可用命令

### Skills

| 命令 | 用途 |
|:---|:---|
| `/gclm` | 智能分流工作流 |
| `/investigate` | 代码库调查 |
| `/tdd` | 测试驱动开发 |
| `/spec` | 规范驱动开发 |
| `/commit` | 智能 Git 提交 |

### gclm-engine

| 命令 | 用途 |
|:---|:---|
| `workflow start` | 创建任务并开始执行 |
| `workflow list` | 列出所有工作流 |
| `task current` | 获取当前待执行阶段 |
| `task complete` | 完成阶段 |
| `task list` | 列出所有任务 |

---

## 目录结构

```
gclm-flow/
├── README.md              # 本文件
├── CLAUDE.md              # Claude Code 项目指导
├── install.sh             # 安装脚本
│
├── gclm-engine/           # Go 引擎
│   ├── main.go           # 入口文件
│   ├── internal/         # 内部实现
│   │   ├── cli/          # CLI 命令
│   │   ├── db/           # 数据库层
│   │   ├── pipeline/     # 流水线解析
│   │   └── service/      # 任务服务
│   ├── pkg/types/        # 共享类型
│   ├── test/             # 测试文件
│   └── Makefile          # 构建脚本
│
├── workflows/             # 工作流 YAML 配置
│   ├── code_simple.yaml  # CODE_SIMPLE 工作流
│   ├── code_complex.yaml # CODE_COMPLEX 工作流
│   ├── document.yaml     # DOCUMENT 工作流
│   └── examples/         # 自定义工作流示例
│
├── agents/                # Agent 定义
│   ├── investigator.md   # 调查 agent
│   ├── architect.md      # 架构设计 agent
│   ├── spec-guide.md     # SpecDD 指导
│   ├── tdd-guide.md      # TDD 指导
│   ├── worker.md         # 代码实现
│   └── code-reviewer.md  # 代码审查
│
├── skills/                # Skill 定义
│   ├── gclm/             # 主工作流 skill
│   └── git-commit/       # Git 提交 skill
│
├── commands/              # 命令定义
├── rules/                 # 工作流规则
├── .github/workflows/     # GitHub Actions
│
└── docs/                  # 文档
```

---

## 核心组件

### gclm-engine (Go 引擎)

- **工作流配置管理** - 基于 YAML 的可扩展工作流定义
- **任务状态管理** - SQLite 持久化，支持暂停/恢复
- **阶段调度** - 依赖解析、拓扑排序、并行执行
- **CLI 接口** - JSON 输出，与 skills 集成

### Skills

- **gclm** - 主工作流执行，调用 Go 引擎创建和执行任务
- **git-commit** - 智能提交，分析项目风格生成提交信息

### Agents

| Agent | 职责 | 模型 |
|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku |
| `architect` | 架构设计、方案权衡 | Opus |
| `spec-guide` | SpecDD 规范文档编写 | Opus |
| `tdd-guide` | TDD 流程指导 | Sonnet |
| `worker` | 执行明确定义的任务 | Sonnet |
| `code-reviewer` | 代码审查 | Sonnet |

---

## 核心特性

- **智能分流**: Go 引擎自动判断任务类型，选择最优工作流
- **llmdoc 优先**: 任何操作前先读取项目文档
- **SpecDD 集成**: 复杂模块先写规范文档
- **TDD 强制**: 测试驱动开发，覆盖率 > 80%
- **状态持久化**: SQLite 存储，支持暂停/恢复
- **可扩展**: YAML 配置工作流，支持自定义

---

## 外部工具

### auggie (推荐)

**语义代码搜索工具**，可提升代码理解效率约 20-30%

```bash
npm install -g @augmentcode/auggie@prerelease
```

---

## 开发

### 构建 gclm-engine

```bash
cd gclm-engine
make build
# 或
make dev    # 快速构建到 ~/.gclm-flow/
```

### 测试

```bash
make test
```

---

## 文件位置

| 组件 | 位置 |
|:---|:---|
| **二进制** | `~/.gclm-flow/gclm-engine` |
| **数据库** | `~/.gclm-flow/gclm-engine.db` |
| **工作流** | `~/.gclm-flow/workflows/` |
| **Agents** | `~/.claude/agents/` |
| **Skills** | `~/.claude/skills/` |

---

## License

MIT

---

<div align="center">

由 **gclm** 精心打造

</div>
