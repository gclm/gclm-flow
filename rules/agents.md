# Agent 编排规则

## 可用 Agents

### gclm-flow 自定义 Agents

| Agent | 职责 | 模型 | 定义文件 |
|:---|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku 4.5 | `agents/investigator.md` |
| `architect` | 架构设计、方案权衡 | Opus 4.5 | `agents/architect.md` |
| `worker` | 执行明确定义的任务 | Sonnet 4.5 | `agents/worker.md` |
| `tdd-guide` | TDD 流程指导 | Sonnet 4.5 | `agents/tdd-guide.md` |
| `spec-guide` | SpecDD 规范文档编写 | Opus 4.5 | `agents/spec-guide.md` |
| `code-reviewer` | 代码审查 | Sonnet 4.5 | `agents/code-reviewer.md` |
| `llmdoc` | LLM 优化文档生成/更新 | Sonnet 4.5 | `agents/llmdoc.md` |

### Claude Code 官方插件

| Agent | 插件名 | 职责 | 何时使用 |
|:---|:---|:---|:---|
| `code-simplifier` | `code-simplifier@claude-plugins-official` | 代码简化重构 | Phase 7 重构优化 |
| `security-guidance` | `security-guidance@claude-plugins-official` | 安全审查 | Phase 7 安全检查 |
| `commit-commands` | `commit-commands@claude-plugins-official` | Git 操作 | Commit/Push/PR |

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

---

## 并行冲突解决

### 输出合并优先级

当多个 Agent 并行输出时，按以下优先级合并：

1. **architect 输出** > **investigator 输出** (设计方案 > 探索发现)
2. **具体文件路径** > **抽象描述**
3. **主 Agent 输出** > **并行 Agent 输出** (最终决策权)

### 矛盾处理流程

```
┌─────────────────┐
│ 并行 Agent 输出  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌──────────────┐
│ 检测矛盾?       │──否─▶│ 直接合并     │
└────────┬────────┘     └──────────────┘
         │是
         ▼
┌─────────────────┐
│ 应用优先级规则   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌──────────────┐
│ 仍有矛盾?       │──否─▶│ 合并结果     │
└────────┬────────┘     └──────────────┘
         │是
         ▼
┌─────────────────┐
│ 展示选项给用户   │ ◀── AskUserQuestion
└─────────────────┘
```

### 具体规则

| 场景 | 处理方式 |
|:---|:---|
| **文件路径冲突** | 具体路径 (含行号) > 抽象路径 |
| **方案冲突** | architect 方案 > investigator 建议 |
| **无法自动解决** | 使用 AskUserQuestion 展示选项 |

### Phase 2 并行探索

3 个 `investigator` 并行时，任务分配：
- Agent 1: 相似功能搜索
- Agent 2: 架构映射
- Agent 3: 代码规范识别

**冲突处理**: 使用不同的 Grep 模式和 Glob 路径，避免文件竞争。

### Phase 4 并行设计

2 个 `architect` + 1 个 `investigator` 并行时：
- Architect 1: 组件设计
- Architect 2: 数据流设计
- Investigator: 依赖分析

**冲突处理**: 各自关注不同层面，最后汇总合并。

### Phase 7 并行审查

`code-simplifier` + `security-guidance` + `code-reviewer` 并行时：
- 各自独立分析
- 输出格式统一化
- 主 Agent 合并建议

**冲突处理**: 优先级：**安全** > **正确性** > **简洁性**

---

## Agent 调用时机表

| 场景 | 触发条件 | Agent | 理由 |
|:---|:---|:---|:---|
| 架构设计 | 涉及 ≥3 个文件或跨模块 | `architect` | 复杂设计需要 Opus |
| Bug 修复 | 单文件，<50 行变更 | `worker` + `tdd-guide` | 简单修复 |
| 代码调查 | 需要理解代码库结构 | `investigator` | 快速调查用 Haiku |
| Spec 文档 | Phase 5 阶段 | `spec-guide` | SpecDD 专用 |
| 测试编写 | Phase 6 阶段 | `tdd-guide` | TDD 指导 |
| 代码审查 | Phase 7 阶段 | `code-simplifier` + `security-guidance` + `code-reviewer` | 并行审查 |
