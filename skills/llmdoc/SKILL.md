---
name: llmdoc
description: LLM 优化的项目文档 - 自动生成或更新 llmdoc
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "Task"]
---

# /llmdoc - LLM 优化的项目文档

**触发时机**: 代码变更后、工作流 summary 阶段完成后、架构变更时、首次使用时

---

## 自动行为

**无需用户确认** - 自动检测并执行：

| 情况 | 行为 |
|:---|:---|
| `llmdoc/` 不存在 | **智能分级生成**（根据项目成熟度） |
| `llmdoc/` 已存在 | 扫描变更并更新文档 |

---

## 执行流程

### 1. 调用 llmdoc Agent

```bash
# 调用 llmdoc Agent 处理文档生成/更新
Task tool: llmdoc agent
```

### 2. Agent 处理

llmdoc Agent 会自动：
1. 检测 llmdoc/ 是否存在
2. 评估项目成熟度（如果不存在）
3. 扫描代码变更（如果已存在）
4. 生成/更新对应文档
5. 验证完整性

### 3. 输出摘要

```markdown
## llmdoc 操作摘要

### 检测结果
- llmdoc 状态: [exists / not_found]
- 项目成熟度: [初级 / 中级 / 成熟]

### 执行操作
- 生成文档: [列表]
- 更新文档: [列表]

### 变更分析
- 变更文件: [列表]
- 影响模块: [列表]
```

---

## 示例

```bash
# 用户输入
/llmdoc

# 自动执行
→ 调用 llmdoc Agent
→ 扫描项目/变更
→ 生成/更新文档
→ 输出摘要
```

---

## 与工作流集成

### 工作流完成后

当工作流的 `summary` 阶段完成后，系统会自动询问：

```
是否需要更新项目文档？
```

- 用户确认 → 调用 llmdoc Agent
- 用户取消 → 跳过

**触发机制**：通过工作流 YAML 中的 `doc_update` 节点（`required: false`）实现可选触发。

---

## Agent 定义

详见: `agents/llmdoc.md`

---

## 相关命令

| 命令 | 用途 |
|:---|:---|
| `/llmdoc` | 生成/更新项目文档 |
| `/gclm` | 智能分流工作流（含文档更新询问） |
