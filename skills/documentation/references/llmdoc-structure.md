# llmdoc Structure Reference

## 目录结构

```
llmdoc/
├── index.md              # 导航和概述
├── overview/             # "这是什么项目？"
│   └── project-overview.md
├── architecture/         # "它是如何工作的？"
│   └── *.md
├── guides/               # "如何做 X？"
│   └── *.md
└── reference/            # "具体细节是什么？"
    ├── conventions/      # 项目约定
    ├── decisions/        # ADR 架构决策记录
    └── *.md
```

## 类别用途

| 类别 | 回答的问题 | 内容类型 |
|------|-----------|----------|
| `overview/` | 这是什么项目？ | 高层上下文、目的、技术栈 |
| `architecture/` | 它是如何工作的？ | 组件关系、系统设计 |
| `guides/` | 如何做 X？ | 分步工作流（最多 5-7 步） |
| `reference/` | 具体细节是什么？ | 约定、数据模型、API 规范 |

## 阅读优先级

1. **首先**阅读 `index.md`
2. **必须**阅读所有 `overview/*.md` 文档
3. 修改相关代码前阅读 `architecture/` 文档
4. 需要步骤指导时查阅 `guides/`
5. 需要规范细节时查看 `reference/`

## 文档约定

- **简洁**：每个文档不超过 150 行
- **引用代码**：使用 `path/file.ext:line` 格式引用
- **命名**：文件名使用 kebab-case，如 `project-overview.md`
- **面向 LLM**：为机器阅读优化

## 初始化内容

### index.md
```markdown
# [项目名称]

## 快速导航

- [项目概述](overview/project-overview.md)
- [开发指南](guides/development.md)
- [架构设计](architecture/)

## 技术栈

- 语言：
- 框架：
- 数据库：

## 快速开始

[基本使用步骤]
```

### overview/project-overview.md
```markdown
# 项目概述

## 简介

[一句话描述]

## 核心功能

- 功能 1
- 功能 2

## 技术架构

[架构概述]

## 相关文档

- [开发指南](../guides/development.md)
- [API 参考](../reference/api.md)
```

### reference/decisions/ADR-001-title.md
```markdown
# ADR-001: [决策标题]

## 状态

提议/已接受/已废弃

## 背景

[描述背景和问题]

## 决策

[描述决策内容]

## 理由

[解释为什么]

## 影响

[描述影响]
```
