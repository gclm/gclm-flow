---
name: gclm-init
description: |
  项目初始化技能。当检测到新项目（无 llmdoc/AGENTS.md）或用户要求初始化、init 时自动触发。
  包含：(1) 检测项目信息 (2) 创建 llmdoc 结构 (3) 生成项目配置
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - init
    - setup
---

# 项目初始化

初始化 Gclm-Flow 工作流，创建项目文档结构。

## 触发条件

- 检测到新项目（无 llmdoc/AGENTS.md）
- 用户说"初始化"、"init"

## 工作流程

### 1. 检测项目信息

调用 `investigator` 代理：
- 自动检测语言和框架
- 识别构建工具
- 检测测试框架

### 2. 创建 llmdoc 结构

```
llmdoc/
├── overview.md          # 项目概述
├── guides/              # 使用指南
│   ├── getting-started.md
│   ├── development.md
│   └── deployment.md
├── architecture/        # 架构文档
│   ├── system-design.md
│   ├── data-model.md
│   └── api-design.md
└── reference/           # 参考资料
    ├── decisions/       # ADR
    ├── conventions/     # 项目约定
    └── patterns/        # 代码模式
```

### 3. 生成项目配置

- 创建 CLAUDE.md（Claude Code）
- 创建 AGENTS.md（Codex CLI）
- 更新 .gitignore

## 语言栈检测

| 检测文件 | 语言 | 模板 |
|----------|------|------|
| `pom.xml` | Java/Maven | Java Spring Boot |
| `build.gradle` | Java/Gradle | Java Spring Boot |
| `requirements.txt` | Python | Python FastAPI |
| `pyproject.toml` | Python | Python FastAPI |
| `go.mod` | Go | Go Gin |
| `Cargo.toml` | Rust | Rust Axum |
| `package.json` | Node.js | 前端 React |

## 选项

- `--force`: 强制覆盖已存在的 llmdoc
- `--minimal`: 只创建最小结构

## 输出

```markdown
# 初始化完成

## 项目信息
- 语言: Java
- 框架: Spring Boot
- 构建工具: Maven
- 测试框架: JUnit 5

## 创建的文件
- llmdoc/overview.md
- llmdoc/guides/getting-started.md
- llmdoc/architecture/system-design.md
- CLAUDE.md

## 下一步
1. 查看 llmdoc/overview.md 了解项目结构
2. 使用 /gclm:plan 规划第一个任务
```
