# Agent 编排规则

## 可用 Agents

| Agent | 职责 | 模型 | 何时使用 |
|:---|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku 4.5 | 代码库调查、快速分析 |
| `architect` | 架构设计、方案权衡 | Opus 4.5 | 架构决策、设计方案 |
| `worker` | 执行明确定义的任务 | Sonnet 4.5 | 代码实现、运行测试 |
| `tdd-guide` | TDD 流程指导 | Sonnet 4.5 | 新功能、Bug 修复 |
| `code-simplifier` | 代码简化重构 | Sonnet 4.5 | Phase 7 重构优化 |
| `security-guidance` | 安全审查 | Sonnet 4.5 | Phase 7 安全检查 |
| `code-reviewer` | 代码审查 | Sonnet 4.5 | 实现后审查 |

---

## Agent 调用时机

### 无需用户提示即可使用

1. **复杂功能请求** → 使用 `architect` agent
2. **代码刚编写/修改** → 使用 `code-simplifier` + `code-reviewer` + `security-guidance` agents（Phase 7）
3. **Bug 修复或新功能** → 使用 `tdd-guide` agent
4. **架构决策** → 使用 `architect` agent
5. **代码库调查** → 使用 `investigator` agent
6. **代码需要简化** → 使用 `code-simplifier` agent
7. **安全审查** → 使用 `security-guidance` agent

---

## 并行任务执行

**始终对独立操作使用并行 Task 执行**

### 示例：并行执行

```markdown
# ✅ 好：并行执行
启动 3 个 agent 并行：
1. Agent 1: 分析 auth.ts 的安全性
2. Agent 2: 审查缓存系统的性能
3. Agent 3: 检查 utils.ts 的类型

# ❌ 坏：不必要的串行
先 agent 1，然后 agent 2，然后 agent 3
```

---

## 多视角分析

对于复杂问题，使用分角色子 agent：
- 事实审查者
- 高级工程师
- 安全专家
- 一致性审查者
- 冗余检查者

---

## Agent 协作模式

### 串行协作

```
investigator (探索)
    ↓
architect (设计)
    ↓
tdd-guide (测试指导)
    ↓
worker (实现)
    ↓
code-reviewer (审查)
```

### 并行协作

```
Phase 2: investigator x3 (并行探索)
    ↓
Phase 4: architect x2 + investigator (并行设计)
    ↓
Phase 7: code-simplifier + security-guidance + code-reviewer (并行重构+安全+审查)
```

---

## 上下文传递

每次 Agent 调用必须包含：

```markdown
## Original User Request
<原始请求>

## Context Pack
- Phase: <阶段名称>
- Decisions: <决策/约束/选择>
- Previous outputs: <之前 Agent 的输出>
- Open questions: <未解决的问题>

## Current Task
<具体任务>

## Acceptance Criteria
<可检查的输出>
```

---

## Agent 输出要求

### investigator
- 简洁报告 (< 150 行)
- 文件引用 (file:line)
- 客观事实
- 无代码粘贴

### architect
- 完整蓝图
- 文件清单
- 组件设计
- 数据流
- 构建序列

### worker
- 最小 diff
- 遵循模式
- 测试通过
- 验证结果

### tdd-guide
- 测试先写
- 失败验证
- 最小实现
- 覆盖率检查

### code-reviewer
- 正确性检查
- 简洁性检查
- 安全性检查
- 可行建议

### code-simplifier
- 保持功能不变
- 提升可读性
- 消除重复代码
- 改进命名和结构
- 优化复杂度

### security-guidance
- OWASP Top 10 检查
- 注入漏洞检测
- 认证授权审查
- 敏感数据处理
- 依赖安全检查

---

## 模型选择策略

| 场景 | 模型 | 原因 |
|:---|:---|:---|
| 快速调查 | Haiku 4.5 | 速度快，成本低 |
| 复杂设计 | Opus 4.5 | 深度思考，高质量 |
| 标准实现 | Sonnet 4.5 | 平衡速度和质量 |
