---
name: dispatching-parallel-agents
description: Use when there are 2 or more independent investigations or implementation tasks that can run in parallel. Trigger on 并行调研, 蜂群模式, 多模块同时, spawn team, agentTeam.
---

# 并行 Agent 调度（蜂群模式）

用 `spawn_team` 工具启动多 agent 协作，时间是约束，并行是手段。

## 使用条件

- 至少 2 个子任务彼此独立
- 每个子任务有独立交付物
- 不会频繁编辑同一文件或依赖同一中间结果

## spawn_team 调用方式

```json
{
  "team_id": "optional-team-name",
  "members": [
    {
      "name": "成员标识",
      "task": "具体任务描述",
      "agent_type": "investigator",
      "worktree": false,
      "background": false
    }
  ]
}
```

agent_type 对应关系：
- `planner` — 任务拆解、计划制定
- `investigator` — 代码调研、技术评估（只读，适合并行）
- `builder` — 代码实现、修改
- `reviewer` — 审查验证（只读）
- `recorder` — 文档沉淀

## 标准团队模式

### 并行调研（最常用）

多个 investigator 同时探索不同方向，汇总后再决策：

```json
{
  "members": [
    {"name": "invest-frontend", "task": "调研前端模块现状", "agent_type": "investigator"},
    {"name": "invest-backend", "task": "调研后端接口现状", "agent_type": "investigator"}
  ]
}
```

### 完整交付团队（顺序分阶段）

调研和实现有依赖关系，分两轮 spawn：

第一轮：先调研
```json
{
  "members": [
    {"name": "researcher", "task": "调研相关代码、约束和现有模式，输出关键发现", "agent_type": "investigator"}
  ]
}
```

第二轮：拿到调研结论后，再实现和验收
```json
{
  "members": [
    {"name": "coder", "task": "基于调研结论实现功能 X，约束：[从第一轮结果填入]", "agent_type": "builder"},
    {"name": "checker", "task": "验收实现结果", "agent_type": "reviewer"}
  ]
}
```

### 独立模块并行实现

多个 builder 在各自 worktree 里独立实现不同模块：

```json
{
  "members": [
    {"name": "build-module-a", "task": "实现模块 A", "agent_type": "builder", "worktree": true},
    {"name": "build-module-b", "task": "实现模块 B", "agent_type": "builder", "worktree": true}
  ]
}
```

## 操作顺序

1. 拆出子任务清单，确认彼此独立。
2. 为每个成员定义明确的 task、输入边界、禁止事项。
3. 调用 `spawn_team`，等待所有成员完成。
4. 汇总结果：改动文件、验证证据、冲突情况。
5. 若发现共享状态或顺序依赖，立即改回串行。

## 不要这样用

- 同一个模块的连续重构（有文件冲突风险）
- 先后强依赖的任务（B 依赖 A 的输出）
- 问题边界还不清晰时
- agent 深度已达上限时（系统报错，改为串行）
