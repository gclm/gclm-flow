# Gclm-Flow

> 全栈开发工作流增强，智能检测语言/框架，统一命令体验

[![npm version](https://badge.fury.io/js/gclm-flow.svg)](https://badge.fury.io/js/gclm-flow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 简介

Gclm-Flow 是一个 Claude Code 工作流增强包，为全栈开发者提供：

- **智能检测**：自动识别项目语言、框架、构建工具
- **统一命令**：12 个命令自动适配不同技术栈
- **记忆系统**：记录错误和解决方案，避免重复犯错
- **6 个代理**：planner、builder、reviewer、investigator、recorder、remember

## 快速开始

### 安装

```bash
npm install -g gclm-flow
```

### 使用

```bash
# 交互式安装
gclm-flow

# 直接安装（使用默认配置）
gclm-flow install -y

# 列出可用组件
gclm-flow list

# 卸载
gclm-flow uninstall
```

## 支持的技术栈

| 语言 | 框架 | 构建工具 |
|------|------|----------|
| Java | Spring Boot | Maven, Gradle |
| Python | Flask, FastAPI | pip, poetry, pipenv |
| Go | Gin, Echo, Fiber | go modules |
| Rust | Axum, Actix, Rocket | cargo |
| 前端 | React, Vue, Angular, Svelte | npm, pnpm, yarn, bun |

## 命令列表

| 命令 | 描述 |
|------|------|
| `/gclm:auto` | 智能命令路由，自动选择合适命令 |
| `/gclm:init` | 初始化项目，创建 llmdoc 文档结构 |
| `/gclm:plan` | 分析需求，制定执行计划 |
| `/gclm:do` | 根据计划执行代码实现 |
| `/gclm:review` | 代码审查，质量检查 |
| `/gclm:test` | 智能运行测试 |
| `/gclm:fix` | 诊断和修复问题 |
| `/gclm:doc` | 文档管理 |
| `/gclm:ask` | 基于项目上下文的问答 |
| `/gclm:learn` | 记忆管理 |
| `/gclm:commit` | 智能生成提交信息 |
| `/gclm:verify` | 对照设计文档验证实现 |

## 代理系统

| 代理 | 职责 |
|------|------|
| **planner** | 任务规划，需求分解 |
| **builder** | 代码实现，修改重构 |
| **reviewer** | 质量检查，测试验证 |
| **investigator** | 上下文调研，技术调研 |
| **recorder** | 文档记录，知识管理 |
| **remember** | 记忆管理，错误记录 |

## Hooks 自动化

| Hook | 功能 |
|------|------|
| Dev server 阻止 | 阻止 tmux 外运行 dev server |
| 自动格式化 | 编辑后自动格式化代码 |
| TypeScript 检查 | 编辑 TS 后类型检查 |
| 调试语句警告 | 检测 console.log 等 |
| 会话管理 | 加载/保存会话状态 |

## 数据存储

所有数据存储在 `~/.gclm-flow/` 目录：

```
~/.gclm-flow/
├── memory/          # 记忆系统
│   ├── errors/      # 错误记忆
│   └── patterns/    # 模式记忆
├── cache/           # 项目检测缓存
├── session.json     # 会话状态
└── config.json      # 配置文件
```

## 项目结构

```
gclm-flow/
├── bin/cli.js           # CLI 入口
├── src/                 # 核心代码
│   ├── index.js         # 主逻辑
│   ├── installer.js     # 安装逻辑
│   ├── prompts.js       # 交互提示
│   ├── utils.js         # 工具函数
│   └── detector.js      # 智能检测
├── agents/              # 6 个代理
├── commands/gclm/       # 12 个命令
├── rules/               # 分层规则
│   ├── core.md          # 核心规则
│   ├── languages/       # 语言规则
│   └── domains/         # 领域规则
├── skills/              # 技能库
│   ├── core/            # 核心工作流
│   ├── patterns/        # 语言模式
│   ├── database/        # 数据库技能
│   ├── devops/          # DevOps 技能
│   └── memory/          # 记忆系统
├── hooks/               # Hooks 配置
├── scripts/hooks/       # Hook 脚本
└── templates/           # 用户配置模板
    ├── user-CLAUDE.md   # 用户级 CLAUDE.md
    ├── statusline.json  # 状态栏配置
    └── mcp-servers.json # MCP 配置
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/gclm/gclm-flow.git
cd gclm-flow

# 安装依赖
npm install

# 本地测试
npm link
gclm-flow list

# 运行测试
npm test
```

## 配置 MCP

Gclm-Flow 包含以下 MCP 服务器配置：

```json
{
  "mcpServers": {
    "auggie": {
      "command": "auggie",
      "args": ["--mcp"]
    },
    "playwright": {
      "command": "npx",
      "args": ["-y", "@playwright/mcp"]
    },
    "exa": {
      "url": "https://mcp.exa.ai/mcp"
    }
  }
}
```

## 相关链接

- [Claude Code](https://claude.ai/code)
- [问题反馈](https://github.com/gclm/gclm-flow/issues)

## License

MIT © gclm
