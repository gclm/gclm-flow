---
name: memory
description: |
  记忆系统技能。当用户要求记住、记录错误、提取模式、查询历史时自动触发。
  包含：(1) 错误记忆 (2) 模式记忆 (3) 知识查询
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - memory
    - learning
---

# 记忆系统

## 数据存储

```
~/.gclm-flow/memory/
├── errors/              # 错误记忆（按语言分类）
│   ├── java.json
│   ├── python.json
│   └── ...
├── patterns/            # 模式记忆（按类型分类）
│   ├── api-design.json
│   └── ...
└── index.json           # 记忆索引
```

## 功能

### 错误记忆
- 记录遇到的错误和解决方案
- 避免重复犯错
- 按语言/框架分类

### 模式记忆
- 提取成功的代码模式
- 促进复用
- 持续改进

### 知识查询
- 查询历史解决方案
- 检索相关模式

## 使用方式

```
/memory save <error|pattern> <content>
/memory query <keyword>
/memory list
```

## Codex 兼容

Codex CLI 使用内置 memories 系统，路径：`~/.codex/memories/`
