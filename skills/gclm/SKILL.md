---
name: gclm
description: "融合工作流 - TDD-First + llmdoc 优先 + 多 Agent 并行。8 阶段流程：llmdoc 读取 → Discovery → Exploration → Clarification → Architecture → TDD Red → TDD Green → Refactor+Doc → Summary"
allowed-tools: ["Bash(${SKILL_DIR}/../scripts/setup-gclm.sh:*)"]
---

# gclm-flow 融合开发工作流 Skill

## 核心哲学

**TDD-First + llmdoc 优先 + 多 Agent 并行**

## 循环初始化 (必需)

当通过 `/gclm <task>` 触发时，**首先**初始化循环状态：

```bash
"${SKILL_DIR}/../scripts/setup-gclm.sh" "<task description>"
```

这会创建 `.claude/gclm.{task_id}.local.md` 包含：
- `active: true`
- `current_phase: 0`
- `max_phases: 8`
- `completion_promise: "<promise>GCLM_WORKFLOW_COMPLETE</promise>"`

## 8 阶段工作流

| 阶段 | 名称 | Agent | 并行 |
|:---|:---|:---|:---:|
| 0 | llmdoc 优先读取 | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 2 | Exploration | `investigator` x3 | 是 |
| 3 | Clarification | 主 Agent + AskUser | - |
| 4 | Architecture | `architect` x2 + `investigator` | 是 |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Doc | `worker` + `code-reviewer` | 是 |
| 8 | Summary | `investigator` | - |

## 硬约束

1. **Phase 0 强制**: 必须优先读取 llmdoc
2. **Phase 3 不可跳过**: 必须澄清所有疑问
3. **Phase 5 TDD 强制**: 必须先写测试
4. **并行优先**: 能并行的任务必须并行执行
5. **状态持久化**: 每个阶段后更新状态文件
6. **选项式编程**: 使用 AskUserQuestion 展示选项
7. **文档更新询问**: Phase 7 必须询问

## 循环状态管理

每个阶段后，更新 `.claude/gclm.{task_id}.local.md` frontmatter：
```yaml
current_phase: <下一阶段编号>
phase_name: "<下一阶段名称>"
```

当所有 8 阶段完成，输出完成信号：
```
<promise>GCLM_WORKFLOW_COMPLETE</promise>
```

提前退出：在状态文件中设置 `active: false`。

## 并行执行示例

### Phase 2: Exploration (3 个并行任务)
```bash
codeagent-wrapper --parallel <<'EOF'
---TASK---
id: p2_similar_features
agent: gclm-investigator
workdir: .
---CONTENT---
Find similar features, trace end-to-end.

---TASK---
id: p2_architecture
agent: gclm-investigator
workdir: .
---CONTENT---
Map architecture for relevant subsystem.

---TASK---
id: p2_conventions
agent: gclm-investigator
workdir: .
---CONTENT---
Identify testing patterns and conventions.
EOF
```

### Phase 4: Architecture (2 个并行方案)
```bash
codeagent-wrapper --parallel <<'EOF'
---TASK---
id: p4_minimal
agent: gclm-architect
workdir: .
---CONTENT---
Propose minimal-change architecture.

---TASK---
id: p4_pragmatic
agent: gclm-architect
workdir: .
---CONTENT---
Propose pragmatic-clean architecture.
EOF
```

## Agent 体系

| Agent | 职责 | 模型 |
|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku 4.5 |
| `architect` | 架构设计、方案权衡 | Opus 4.5 |
| `worker` | 执行明确定义的任务 | Sonnet 4.5 |
| `tdd-guide` | TDD 流程指导 | Sonnet 4.5 |
| `code-reviewer` | 代码审查 | Sonnet 4.5 |

## 上下文包模板

```text
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

## Stop Hook

注册 Stop Hook 后：
1. 创建 `.claude/gclm.{task_id}.local.md` 状态文件
2. 每个阶段后更新 `current_phase`
3. Stop hook 检查状态，未完成时阻止退出
4. 完成时输出 `<promise>GCLM_WORKFLOW_COMPLETE</promise>`

手动退出：在状态文件中设置 `active` 为 `false`。
