# Doc Workflow

管理 llmdoc 文档系统的工作流技能。

## 何时使用

- 初始化项目时创建 llmdoc 结构
- 更新项目文档
- 生成 API 文档
- 记录架构决策

## llmdoc 结构

参考 [llmdoc-structure.md](references/llmdoc-structure.md) 了解完整结构。

## 工作流程

### 1. 初始化文档

```
/gclm:init
```

1. 检测项目技术栈
2. 创建 llmdoc 目录结构
3. 生成项目概述
4. 创建初始指南

### 2. 更新文档

```
/gclm:doc update
```

1. 分析代码变更
2. 更新受影响的文档
3. 保持文档与代码同步

### 3. 记录决策

```
/gclm:doc decision "决策标题"
```

在 `llmdoc/reference/decisions/` 创建 ADR 文档。

## 文档规范

参考 [doc-conventions.md](references/doc-conventions.md) 了解写作规范。

## 相关命令

- `/gclm:init` - 初始化项目
- `/gclm:doc` - 文档管理
- `/gclm:ask` - 基于文档回答问题
