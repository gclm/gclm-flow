# Agent 体系

## Agent 概览

gclm-flow 使用 7 个自定义 Agent + 2 个官方插件 Agent，覆盖开发生命周期的各个阶段。

```
┌─────────────────────────────────────────────────────────────┐
│                      gclm-flow Agent 架构                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   自定义 Agents (7个)                                        │
│   ┌──────────────┬──────────────┬──────────────┐            │
│   │ investigator │   architect  │  spec-guide  │            │
│   │  Haiku 4.5   │  Opus 4.5    │  Opus 4.5    │            │
│   └──────────────┴──────────────┴──────────────┘            │
│   ┌──────────────┬──────────────┬──────────────┐            │
│   │  tdd-guide   │    worker    │code-reviewer │            │
│   │ Sonnet 4.5   │ Sonnet 4.5   │ Sonnet 4.5   │            │
│   └──────────────┴──────────────┴──────────────┘            │
│   ┌──────────────────────────────────────────────┐          │
│   │            recorder                           │          │
│   │         Sonnet 4.5                            │          │
│   └──────────────────────────────────────────────┘          │
│                                                              │
│   官方插件 Agents (2个)                                      │
│   ┌──────────────────┬──────────────────────────┐            │
│   │ code-simplifier  │  security-guidance       │            │
│   └──────────────────┴──────────────────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

---

## 自定义 Agents

### investigator

| 属性 | 值 |
|:---|:---|
| **模型** | Haiku 4.5 |
| **职责** | 代码库调查、分析、总结 |
| **使用阶段** | Phase 1, 2, 9 |
| **文件** | `agents/investigator.md` |

**核心能力**:
- 语义代码搜索 (auggie 优先)
- 代码库结构映射
- 相似功能发现
- 代码规范识别

**输出特点**: 简洁报告 (< 150 行)，客观事实，无代码粘贴

---

### architect

| 属性 | 值 |
|:---|:---|
| **模型** | Opus 4.5 |
| **职责** | 架构设计、方案权衡 |
| **使用阶段** | Phase 4 |
| **文件** | `agents/architect.md` |

**核心能力**:
- 组件关系图设计
- 技术选型分析
- 目录结构规划
- 数据流设计

**输出特点**: 完整蓝图，文件清单，组件设计，构建序列

---

### spec-guide

| 属性 | 值 |
|:---|:---|
| **模型** | Opus 4.5 |
| **职责** | SpecDD 规范文档编写 |
| **使用阶段** | Phase 5 |
| **文件** | `agents/spec-guide.md` |

**核心能力**:
- 编写详细规范文档
- API 接口定义
- 数据结构设计
- 测试策略规划

**输出文件**: `.claude/specs/{feature-name}.md`

---

### tdd-guide

| 属性 | 值 |
|:---|:---|
| **模型** | Sonnet 4.5 |
| **职责** | TDD 流程指导 |
| **使用阶段** | Phase 6 |
| **文件** | `agents/tdd-guide.md` |

**核心能力**:
- 测试驱动开发指导
- Red-Green-Refactor 循环
- 测试覆盖率检查
- 边缘情况识别

**绝对规则**: 绝不一次性生成代码和测试

---

### worker

| 属性 | 值 |
|:---|:---|
| **模型** | Sonnet 4.5 |
| **职责** | 执行明确定义的任务 |
| **使用阶段** | Phase 7 |
| **文件** | `agents/worker.md` |

**核心能力**:
- 最小 diff 实现
- 遵循现有代码模式
- 测试通过验证
- 结果确认

**原则**: 不添加未要求的功能，不过度设计

---

### code-reviewer

| 属性 | 值 |
|:---|:---|
| **模型** | Sonnet 4.5 |
| **职责** | 代码审查 |
| **使用阶段** | Phase 8 |
| **文件** | `agents/code-reviewer.md` |

**核心能力**:
- 正确性检查
- 简洁性检查
- 安全性检查
- 可行建议

---

### recorder

| 属性 | 值 |
|:---|:---|
| **模型** | Sonnet 4.5 |
| **职责** | 文档记录与更新 |
| **使用阶段** | Phase 8 (按需) |
| **文件** | `agents/recorder.md` |

**核心能力**:
- 追踪代码变更影响范围
- 更新 llmdoc 文档
- 验证文档完整性
- 保持文档与代码同步

**使用场景**: Phase 8 完成后，系统会询问是否更新项目文档

---

## 官方插件 Agents

### code-simplifier

| 属性 | 值 |
|:---|:---|
| **插件** | `code-simplifier@claude-plugins-official` |
| **职责** | 代码简化重构 |
| **使用阶段** | Phase 8 |

**核心能力**:
- 消除重复代码
- 改进命名和结构
- 优化复杂度
- 保持功能不变

---

### security-guidance

| 属性 | 值 |
|:---|:---|
| **插件** | `security-guidance@claude-plugins-official` |
| **职责** | 安全审查 |
| **使用阶段** | Phase 8 |

**核心能力**:
- OWASP Top 10 检查
- 注入漏洞检测
- 认证授权审查
- 敏感数据处理

---

## Agent 协作模式

### 串行协作

```
investigator (探索)
    ↓
architect (设计)
    ↓
spec-guide (规范)
    ↓
tdd-guide (测试指导)
    ↓
worker (实现)
    ↓
code-reviewer (审查)
```

### 并行协作

**Phase 2**: investigator x3 (并行探索)
**Phase 4**: architect x2 + investigator (并行设计)
**Phase 8**: code-simplifier + security-guidance + code-reviewer (并行重构+安全+审查)

---

## 上下文传递格式

每次 Agent 调用包含：

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
