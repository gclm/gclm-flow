# gclm-flow

> 智能工作流引擎 - 基于 Go 的任务调度系统，通过 YAML 配置工作流，协调多 Agent 协作完成任务

---

## 快速导航

| 我想... | 立即跳转 |
|:---|:---|
| **快速安装** | [安装指南](#安装) |
| **基本使用** | [使用方法](#使用方法) |
| **工作流介绍** | [工作流概览](#工作流概览) |
| **可用命令** | [命令参考](#命令参考) |
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
2. 安装 agents、skills 到 `~/.claude/`
3. 同步工作流 YAML 配置文件
4. 初始化 SQLite 数据库

---

## 使用方法

### 基本使用

```bash
# 智能工作流（自动判断任务类型）
/gclm 实现用户登录功能

# 代码库调查
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

# 创建任务
gclm-engine task create "修复登录页面 bug" --workflow fix

# 查看当前待执行阶段
gclm-engine task current <task-id>

# 完成阶段
gclm-engine task complete <task-id> <phase-id> --output "执行结果"

# 列出所有任务
gclm-engine task list

# 启动 Web UI 服务器
gclm-engine serve --port 9988
```

---

## 工作流概览

### 内置工作流

| 名称 | 类型 | 适用场景 |
|:---|:---:|:---|
| **analyze** | analyze | 代码分析、问题诊断、性能评估、安全审计 |
| **docs** | docs | 文档编写、设计方案、需求分析 |
| **feat** | feat | 新功能开发、模块开发、跨文件重构 |
| **fix** | fix | Bug 修复、小修改、单文件变更 |

### 工作流配置

工作流由 **YAML 配置文件**定义，位于 `gclm-engine/workflows/` 目录：

```
workflows/
├── analyze.yaml   # 代码分析工作流
├── docs.yaml      # 文档编写工作流
├── feat.yaml      # 复杂功能开发工作流
└── fix.yaml       # Bug 修复工作流
```

每个工作流包含多个阶段（Phase），每个阶段指定 Agent、模型和依赖关系。

---

## 命令参考

### workflow 命令

| 命令 | 说明 |
|:---|:---|
| `workflow list` | 列出所有工作流 |
| `workflow info <name>` | 显示工作流详情 |
| `workflow validate <file>` | 验证 YAML 配置 |
| `workflow install <file>` | 安装工作流 |
| `workflow sync [file]` | 同步 YAML 到数据库 |
| `workflow export <name>` | 导出工作流配置 |

### task 命令

| 命令 | 说明 |
|:---|:---|
| `task create <prompt> --workflow <name>` | 创建任务（使用指定工作流） |
| `task get <task-id>` | 获取任务详情 |
| `task list` | 列出所有任务 |
| `task current <task-id>` | 获取当前待执行阶段 |
| `task plan <task-id>` | 获取执行计划 |
| `task complete <task-id> <phase-id> --output <text>` | 完成阶段 |
| `task fail <task-id> <phase-id> --error <msg>` | 标记阶段失败 |
| `task phases <task-id>` | 显示任务阶段 |
| `task events <task-id>` | 显示任务事件 |
| `task pause/resume/cancel <task-id>` | 暂停/恢复/取消任务 |

### 其他命令

| 命令 | 说明 |
|:---|:---|
| `init` | 初始化配置 |
| `serve --port <port>` | 启动 HTTP API + Web UI 服务器 |
| `version` | 显示版本信息 |

### 全局标志

| 标志 | 说明 |
|:---|:---|
| `--json, -j` | 输出 JSON 格式（便于脚本解析） |

---

## 目录结构

```
gclm-flow/
├── README.md              # 本文件
├── CLAUDE.example.md      # Claude Code 项目指导
├── install.sh             # 安装脚本
│
├── gclm-engine/           # Go 引擎
│   ├── main.go           # 入口文件
│   ├── internal/         # 内部实现
│   │   ├── api/          # HTTP API + WebSocket
│   │   ├── assets/       # 嵌入的 Web UI 文件
│   │   ├── cli/          # CLI 命令
│   │   ├── db/           # 数据库层
│   │   ├── domain/       # 领域接口
│   │   ├── errors/       # 错误处理
│   │   ├── logger/       # 日志系统
│   │   ├── repository/   # 数据仓储
│   │   ├── service/      # 业务服务
│   │   └── workflow/     # 工作流解析器
│   ├── migrations/       # 数据库迁移
│   ├── pkg/types/        # 共享类型
│   ├── test/             # 测试文件
│   ├── web/              # Web UI 静态文件
│   └── workflows/        # 工作流 YAML 配置
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
│   └── gclm/             # 主工作流 skill
│
├── commands/              # 命令定义
├── rules/                 # 工作流规则
├── llmdoc/                # LLM 优化的项目文档
└── docs/                  # 文档
```

---

## 核心组件

### gclm-engine (Go 引擎)

| 功能 | 说明 |
|:---|:---|
| **工作流配置** | 基于 YAML 的可扩展工作流定义 |
| **状态管理** | SQLite 持久化，支持暂停/恢复 |
| **阶段调度** | 依赖解析、拓扑排序、并行执行 |
| **CLI 接口** | JSON 输出，与 skills 集成 |
| **HTTP API** | RESTful API + WebSocket 实时更新 |
| **Web UI** | 可视化任务和阶段状态 |

### Skills

| Skill | 用途 |
|:---|:---|
| `/gclm` | 智能分流工作流，自动判断任务类型 |
| `/investigate` | 代码库调查 |
| `/tdd` | 测试驱动开发 |
| `/spec` | 规范驱动开发 |
| `/commit` | 智能 Git 提交 |

### Agents

| Agent | 职责 | 模型 |
|:---|:---|:---:|
| `investigator` | 探索、分析、总结 | Haiku |
| `architect` | 架构设计、方案权衡 | Opus |
| `spec-guide` | SpecDD 规范文档 | Opus |
| `tdd-guide` | TDD 流程指导 | Sonnet |
| `worker` | 执行明确定义的任务 | Sonnet |
| `code-reviewer` | 代码审查 | Sonnet |

---

## 核心特性

- **智能分流**: 自动判断任务类型，选择最优工作流
- **状态持久化**: SQLite 存储，支持暂停/恢复
- **YAML 可配置**: 工作流完全可定制
- **并行执行**: 支持阶段并行执行
- **JSON 输出**: 便于脚本集成
- **Web UI**: 可视化任务状态和进度

---

## Web UI

启动服务器后，访问 `http://localhost:9988` 可以看到：

- 📊 **仪表板** - 任务统计概览
- 📋 **任务管理** - 创建、查看、管理任务
- ⚙️ **工作流管理** - 查看工作流配置、图示、YAML

```bash
gclm-engine serve --port 9988
```

---

## 开发

### 构建

```bash
cd gclm-engine
make build    # 生产构建
make dev      # 快速构建到 ~/.gclm-flow/
```

### 测试

```bash
make test
```

### 交叉编译

```bash
make build-linux-arm64
make build-darwin-arm64
make build-windows
```

---

## 文件位置

| 组件 | 位置 |
|:---|:---|
| **二进制** | `~/.gclm-flow/gclm-engine` |
| **数据库** | `~/.gclm-flow/gclm-engine.db` |
| **工作流配置** | `~/.gclm-flow/workflows/` |
| **日志文件** | `~/.gclm-flow/gclm-engine.log` |
| **Agents** | `~/.claude/agents/` |
| **Skills** | `~/.claude/skills/` |

---

## 技术栈

| 组件 | 技术 |
|:---|:---|
| **引擎** | Go 1.23+ |
| **数据库** | SQLite (goose migrations) |
| **CLI** | cobra |
| **API** | gorilla/mux + WebSocket |
| **前端** | 嵌入式静态文件 |

---

## License

MIT

---

<div align="center">

由 **gclm** 精心打造

</div>
