# gclm-flow 工作流详细文档

## 融合来源

### myclaude/do - 7 阶段框架

| 特点 | 描述 |
|:---|:---|
| 结构化流程 | 7 个清晰定义的阶段 |
| 状态持久化 | 状态文件支持中途恢复 |
| 强制澄清 | Phase 3 不可跳过 |
| 审批门 | Phase 5 需要用户审批 |
| 并行执行 | 支持多 agent 并行 |

### cc-plugin/Context Floor - llmdoc 解决方案

| 特点 | 描述 |
|:---|:---|
| llmdoc 系统 | LLM 优化的文档结构 |
| SubAgent RAG | 文档优先的调研方法 |
| 文档优先读取 | Phase 0 强制执行 |
| 自动维护提示 | 完成后询问是否更新文档 |

### everything-claude-code/TDD - 实现方法

| 特点 | 描述 |
|:---|:---|
| Red-Green-Refactor | 标准的 TDD 循环 |
| 80% 覆盖率 | 强制测试覆盖率要求 |
| 测试优先 | 绝不一次性生成代码和测试 |
| 测试先失败 | Phase 5 必须先写失败测试 |

---

## 8 阶段详细流程

### Phase 0: llmdoc 优先读取

**目标**: 建立上下文理解

**步骤**:
1. 检查 `llmdoc/` 是否存在
2. 读取 `llmdoc/index.md`
3. 读取 `llmdoc/overview/*.md` (全部)
4. 根据任务读取 `llmdoc/architecture/*.md`

**输出**: 上下文摘要

**强制**: 此阶段不可跳过

---

### Phase 1: Discovery - 理解需求

**Agent**: `investigator`

**活动**:
1. 使用 `AskUserQuestion` 了解用户可见行为、范围、约束
2. 调用 `investigator` 生成需求清单和澄清问题

**输出**:
- Requirements (需求)
- Non-goals (非目标)
- Risks (风险)
- Acceptance Criteria (验收标准)
- Questions (<= 10 个问题)

---

### Phase 2: Exploration - 探索代码库

**并行执行 3 个 `investigator`**

| 任务 | 描述 | 输出 |
|:---|:---|:---|
| 相似功能 | 查找 1-3 个相似功能 | 关键文件、调用流程、扩展点 |
| 架构映射 | 映射相关子系统 | 模块图 + 5-10 个关键文件 |
| 代码规范 | 识别测试模式、规范 | 测试命令 + 文件位置 |

**并行执行**: 必须在单个响应中使用多个 Task 调用

---

### Phase 3: Clarification - 澄清疑问

**强制阶段，不可跳过**

**活动**:
1. 汇总 Phase 1 和 Phase 2 输出
2. 生成优先级排序的问题列表
3. 使用 `AskUserQuestion` 逐一确认

**约束**: 不回答完不进入下一阶段

---

### Phase 4: Architecture - 设计方案

**并行执行**: 2 个 `architect` + 1 个 `investigator`

| Agent | 任务 |
|:---|:---|
| architect (minimal) | 最小改动方案 - 复用现有抽象 |
| architect (pragmatic) | 务实整洁方案 - 引入测试友好接缝 |
| investigator | 测试策略分析 |

**活动**:
1. 使用 `AskUserQuestion` 让用户选择方案
2. 显式审批门: "Approve starting implementation?"

**输出**:
- 文件清单（创建/修改）
- 组件设计
- 数据流
- 构建序列

---

### Phase 5: TDD Red - 编写测试

**Agent**: `tdd-guide`

**TDD 约束**:
- 绝不一次性生成代码和测试
- 先写测试，后写实现
- 测试必须先失败
- 覆盖率目标: > 80%

**流程**:
1. 定义接口
2. 编写测试
3. 运行测试确认失败

**测试应包含**:
- 快乐路径
- 边缘情况
- 错误处理
- 边界值

---

### Phase 6: TDD Green - 编写实现

**Agent**: `worker`

**约束**:
- diff 最小化
- 遵循现有代码模式
- 运行最窄范围的相关测试

**流程**:
1. 编写实现
2. 运行测试
3. 检查覆盖

**验证**:
- 所有测试通过
- 覆盖率 > 80%
- 代码风格一致

---

### Phase 7: Refactor + Doc - 重构与更新文档

**并行执行**:

| Agent | 任务 |
|:---|:---|
| worker | 代码重构 - 优化结构、消除重复 |
| code-reviewer | 代码审查 - 正确性 + 简洁性 |

**重构原则**:
- 保持测试绿色
- 消除重复
- 改进命名
- 优化性能

**文档更新询问**:
```
AskUserQuestion: "是否使用 recorder agent 更新项目文档？"
```

---

### Phase 8: Summary - 完成总结

**Agent**: `investigator`

**输出**:
- 完成的工作内容
- 关键决策和取舍
- 修改的文件路径
- 验证命令
- 后续工作建议

**完成信号**: `<promise>GCLM_WORKFLOW_COMPLETE</promise>`

---

## 状态管理

### 状态文件

```
.claude/gclm.{task_id}.local.md
```

### 状态文件格式

```yaml
---
active: true
current_phase: 0
phase_name: "llmdoc Reading"
max_phases: 8
completion_promise: "<promise>GCLM_WORKFLOW_COMPLETE</promise>"

phases:
  - phase: 0
    name: "llmdoc Reading"
    status: "in_progress"
    started_at: "2026-01-26T15:00:00Z"
---
```

### 状态更新

每个阶段完成后更新:
```yaml
current_phase: <下一阶段编号>
phase_name: "<下一阶段名称>"
```

---

## 并行执行模式

### 必须并行的阶段

- **Phase 2**: 3 个 investigator
- **Phase 4**: 2 个 architect + 1 个 investigator
- **Phase 7**: worker + code-reviewer

### 并行执行格式

```javascript
// 单个响应中的多个 Task 调用
[
  Task({ subagent_type: "investigator", description: "任务1", ... }),
  Task({ subagent_type: "investigator", description: "任务2", ... }),
  Task({ subagent_type: "investigator", description: "任务3", ... })
]
```

---

## Stop Hook

### 位置

`~/.claude/hooks/stop/gclm-loop-hook.sh`

### 行为

1. 检查 `.claude/gclm.*.local.md` 状态文件
2. 如果 `active: true` 且未完成，阻止退出
3. 显示当前阶段和警告
4. 提供强制退出方法

### 强制退出

```bash
sed -i.bak 's/^active: true/active: false/' .claude/gclm.*.local.md
```

---

## 上下文包模板

```markdown
## Original User Request
<verbatim request>

## Context Pack
- Phase: <0-8 name>
- Decisions: <requirements/constraints/choices>
- Investigator output: <paste or "None">
- Architect output: <paste or "None">
- Worker output: <paste or "None">
- Code-reviewer output: <paste or "None">
- Tdd-guide output: <paste or "None">
- Open questions: <list or "None">

## Current Task
<specific task>

## Acceptance Criteria
<checkable outputs>
```
