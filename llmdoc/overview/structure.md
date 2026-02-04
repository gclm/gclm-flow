# 目录结构

## 项目结构

```
gclm-flow/
├── README.md                  # 项目文档
├── install.sh                 # 全局安装脚本
├── CLAUDE.example.md          # 核心配置示例
├── settings.example.json      # MCP 配置示例
│
├── agents/                    # 自定义 Agent 定义
│   ├── investigator.md        # 代码库调查 Agent
│   ├── architect.md           # 架构设计 Agent
│   ├── worker.md              # 任务执行 Agent
│   ├── tdd-guide.md           # TDD 指导 Agent
│   ├── spec-guide.md          # SpecDD 指导 Agent
│   └── code-reviewer.md       # 代码审查 Agent
│
├── commands/                  # Claude Code 命令定义
│   ├── gclm.md                # 智能分流工作流命令
│   ├── investigate.md         # 代码库调查命令
│   ├── tdd.md                 # 测试驱动开发命令
│   ├── spec.md                # 规范驱动开发命令
│   └── llmdoc.md              # 文档生成命令
│
├── skills/                    # 核心工作流 Skill
│   └── gclm/
│       ├── SKILL.md           # Skill 定义和逻辑
│       └── setup-gclm.sh      # 状态初始化脚本
│
├── rules/                     # 工作流规则文档
│   ├── agents.md              # Agent 编排规则
│   ├── llmdoc.md              # llmdoc 文档规则
│   ├── phases.md              # 阶段流程规则
│   ├── spec.md                # SpecDD 规范
│   └── tdd.md                 # TDD 规范
│
├── hooks/                     # Claude Code Hooks
│   ├── notify.sh              # 通知 Hook
│   └── stop-gclm-loop.sh      # 停止 Hook
│
├── .claude-plugin/            # 插件元数据
│   ├── plugin.json            # 插件配置
│   └── marketplace.json       # 市场配置
│
└── llmdoc/                    # LLM 优化文档
    ├── index.md               # 文档索引
    ├── overview/              # 项目概览
    ├── architecture/          # 架构设计
    ├── guides/                # 使用指南
    └── reference/             # 参考文档
```

---

## 目录说明

### `agents/`

自定义 Agent 定义文件，每个 Agent 包含：
- 职责描述
- 模型选择
- 输入输出格式
- 使用场景

### `commands/`

Claude Code 命令定义，用户可通过 `/command` 语法调用：
- `/gclm` - 智能分流工作流
- `/investigate` - 代码库调查
- `/tdd` - 测试驱动开发
- `/spec` - 规范驱动开发
- `/llmdoc` - 文档生成/更新

### `skills/`

核心工作流 Skill，包含：
- `SKILL.md` - Skill 定义和完整工作流逻辑
- `setup-gclm.sh` - 状态文件初始化脚本

### `rules/`

工作流规则文档，定义：
- Agent 编排规则
- llmdoc 文档规则
- 阶段流程规则
- SpecDD/TDD 规范

### `hooks/`

Claude Code Hooks：
- `notify.sh` - 通知用户工作流状态
- `stop-gclm-loop.sh` - 停止循环时清理状态

### `.claude-plugin/`

插件元数据，用于 Claude Code 插件市场：
- `plugin.json` - 插件配置和依赖
- `marketplace.json` - 市场展示信息

### `llmdoc/`

LLM 优化的项目文档：
- `index.md` - 导航入口
- `overview/` - 项目概览
- `architecture/` - 架构设计
- `guides/` - 使用指南
- `reference/` - 参考文档
