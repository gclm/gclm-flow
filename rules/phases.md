# gclm-flow 阶段规则

## 8 阶段工作流详细规则

### Phase 0: llmdoc 优先读取 (NON-NEGOTIABLE)

**目标**: 在任何代码操作前建立上下文理解

**步骤**:
1. 检查 `llmdoc/` 是否存在
2. **如果存在**:
   - 读取 `llmdoc/index.md` 获取导航
   - 读取 `llmdoc/overview/*.md` 全部文档
   - 根据任务读取相关 `llmdoc/architecture/*.md`
3. **如果不存在**:
   - 使用 `investigator` agent 扫描代码库
   - 自动生成 `llmdoc/index.md`
   - 自动生成 `llmdoc/overview/` 基础文档
   - 然后读取生成的文档

**输出**: 上下文摘要（关键文件、模块依赖、设计模式）

**强制**: 此阶段不可跳过
**自动化**: llmdoc 不存在时自动生成，无需用户确认

---

### Phase 1: Discovery - 理解需求

**Agent**: `investigator`

**输出**:
- Requirements (需求)
- Non-goals (非目标)
- Risks (风险)
- Acceptance Criteria (验收标准)
- Questions (澄清问题 <= 10 个)

**关键问题**:
- 用户可见行为是什么？
- 范围和边界在哪里？
- 有哪些约束和限制？
- 成功的标准是什么？

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

### Phase 3: Clarification - 澄清疑问 (强制)

**不可跳过的阶段**

1. 汇总 Phase 1 和 Phase 2 输出
2. 生成优先级排序的问题列表
3. 使用 `AskUserQuestion` 逐一确认

**约束**: 不回答完不进入下一阶段

**典型问题**:
- 技术选型确认
- 数据结构确认
- 错误处理策略
- 性能要求
- 兼容性要求

---

### Phase 4: Architecture - 设计方案

**并行执行**: 2 个 `architect` + 1 个 `investigator`

| Agent | 任务 |
|:---|:---|
| architect (minimal) | 最小改动方案 - 复用现有抽象 |
| architect (pragmatic) | 务实整洁方案 - 引入测试友好接缝 |
| investigator | 测试策略分析 |

使用 `AskUserQuestion` 选择方案

**显式审批门**: "Approve starting implementation?"

**输出应包含**:
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
- 测试必须先失败 (Red)
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

### 状态文件位置

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
  - phase: 1
    name: "Discovery"
    status: "pending"
---
```

### 状态更新

**自动化**: 状态文件更新**自动进行**，无需用户确认

每个阶段完成后自动更新：
```yaml
current_phase: <下一阶段编号>
phase_name: "<下一阶段名称>"
```

**自动化原因**:
- 状态文件是内部元数据，不是代码
- 更新是确定性的（阶段完成 → 状态更新）
- 不影响代码质量或安全性

**仍需授权的场景**:
- Phase 4: Architecture 设计方案审批
- Phase 7: 文档更新询问

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

**重要**: 必须在单个响应中完成所有 Task 调用以实现并行。

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
