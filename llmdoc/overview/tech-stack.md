# 技术栈

## 核心技术

| 类别 | 技术 | 版本 | 用途 |
|:---|:---|:---|:---|
| **编程语言** | Go | 1.21+ | gclm-engine 核心 |
| **数据库** | SQLite | 3.x | 状态存储 (WAL 模式) |
| **迁移工具** | Goose | 3.x | 数据库版本管理 |
| **数据库驱动** | mattn/go-sqlite3 | - | SQLite CGO 驱动 |
| **CLI 框架** | spf13/cobra | - | 命令行接口 |
| **YAML 解析** | gopkg.in/yaml.v3 | - | 配置文件解析 |
| **资源嵌入** | embed (Go 1.16+) | - | 零依赖部署 |
| **脚本语言** | Bash/Zsh | - | 安装脚本、Hooks |
| **配置格式** | Markdown | - | Agent 定义、规则文档 |
| **配置格式** | YAML | - | 工作流定义、前置元数据 |
| **配置格式** | JSON | - | 插件配置、MCP 设置 |
| **插件系统** | Claude Code Plugin API | - | 插件开发和注册 |
| **通信协议** | MCP (Model Context Protocol) | - | 外部工具通信 |

---

## Go 依赖

### 核心依赖

| 包 | 版本 | 用途 |
|:---|:---|:---|
| `github.com/spf13/cobra` | latest | CLI 框架 |
| `github.com/mattn/go-sqlite3` | latest | SQLite 驱动 (CGO) |
| `github.com/pressly/goose/v3` | v3.17.0+ | 数据库迁移 |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML 解析 |

### 构建依赖

- **Go 版本**: 1.21 或更高
- **CGO**: 必需 (go-sqlite3)
- **交叉编译**: 支持 Linux/macOS/Windows (ARM64/AMD64)

---

## AI 模型

| 模型 | 用途 | 使用场景 |
|:---|:---|:---|
| **Claude Opus 4.5** | 高质量设计 | 架构设计 (architect)、规范文档 (spec-guide) |
| **Claude Sonnet 4.5** | 标准实现 | 任务执行 (worker)、TDD 指导 (tdd-guide)、代码审查 (code-reviewer) |
| **Claude Haiku 4.5** | 快速调查 | 代码库调查 (investigator) |

---

## 外部工具

### auggie (推荐)

**描述**: 语义代码搜索工具

**安装**:
```bash
npm install -g @augmentcode/auggie@prerelease
```

**用途**:
- 自然语言代码搜索
- 代码上下文增强
- 语义代码理解

**状态**: 可选但强烈推荐

### 官方插件 Agents

| 插件 | 安装命令 | 用途 |
|:---|:---|:---|
| `code-simplifier` | 内置或通过插件市场 | 代码简化重构 |
| `security-guidance` | 内置或通过插件市场 | 安全审查 |

---

## Claude Code Hooks

支持的 Hook 类型：

| Hook 类型 | 触发时机 | 当前使用 |
|:---|:---|:---|
| `PermissionRequest` | 需要权限确认时 | `notify.sh` |
| `Notification` | 有待处理的操作时 | `notify.sh` |
| `Stop` | 用户尝试退出时 | (未注册) |

---

## 目录权限

安装后的目录结构权限：

```
~/.claude/
├── hooks/
│   ├── notify.sh              (755 可执行)
│   └── stop-gclm-loop.sh      (755 可执行)
├── agents/                    (644)
│   ├── investigator.md
│   ├── architect.md
│   ├── worker.md
│   ├── tdd-guide.md
│   ├── spec-guide.md
│   └── code-reviewer.md
└── rules/                     (644)
    ├── agents.md
    ├── llmdoc.md
    ├── phases.md
    ├── spec.md
    └── tdd.md
```
