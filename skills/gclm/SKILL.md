---
name: gclm
description: "智能分流工作流 - SpecDD + TDD + llmdoc 优先 + ace-tool + 多 Agent 并行。自动判断任务类型：简单任务走 TDD，复杂任务走 SpecDD (Architecture + Spec + TDD)"
allowed-tools: ["Bash(${SKILL_DIR}/../scripts/setup-gclm.sh:*)"]
---

# gclm-flow 智能分流工作流 Skill

## 核心哲学

**SpecDD + TDD + llmdoc 优先 + auggie + 多 Agent 并行 + 智能分流**

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

## 智能分流工作流

### 工作流程图

```
                    ┌─────────────────────────────────────┐
                    │         Phase 1: Discovery           │
                    │    (分析需求、预估规模、复杂度)       │
                    └──────────────────┬──────────────────┘
                                       │
                        ┌──────────────┴──────────────┐
                        │    智能分类 (自动判断)        │
                        └──────────────┬──────────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ↓                        ↓                        ↓
    ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
    │   简单任务       │      │   中等任务       │      │   复杂任务       │
    │ (Bug/小修改)     │      │ (需要确认)       │      │ (新功能/模块)    │
    └────────┬────────┘      └────────┬────────┘      └────────┬────────┘
             │                        │                        │
             │                        ▼                        │
             │              ┌──────────────────┐              │
             │              │ 询问用户确认    │              │
             │              │ 走哪个流程？     │              │
             │              └────────┬─────────┘              │
             │                       │                        │
             │           ┌───────────┴───────────┐           │
             │           │                       │           │
             │           ↓                       ↓           │
             │    ┌─────────────┐        ┌─────────────┐   │
             │    │   简单流程   │        │   完整流程   │   │
             │    └──────┬──────┘        └──────┬──────┘   │
             │           │                      │          │
             └───────────┼──────────────────────┼──────────┘
                         │                      │
                         ▼                      ▼
              ┌─────────────────────┐   ┌─────────────────────┐
              │  Phase 3: Clarify    │   │ Phase 2: Explore   │
              │  (确认问题)          │   │ (并行探索 x3)       │
              └──────────┬──────────┘   └──────────┬──────────┘
                         │                      │
                         ▼                      ▼
              ┌─────────────────────┐   ┌─────────────────────┐
              │  Phase 5: TDD Red    │   │Phase 3: Clarify     │
              │  (写测试)            │   │(澄清疑问)           │
              └──────────┬──────────┘   └──────────┬──────────┘
                         │                      │
                         ▼                      ▼
              ┌─────────────────────┐   ┌─────────────────────┐
              │ Phase 6: TDD Green   │   │Phase 4: Architecture│
              │ (写实现)             │   │(架构设计 x2 + inv)  │
              └──────────┬──────────┘   └──────────┬──────────┘
                         │                      │
                         ▼                      ▼
              ┌─────────────────────┐   ┌─────────────────────┐
              │ Phase 7: Refactor    │   │Phase 4.5: Spec     │
              │ (重构+审查)          │   │(编写规范文档)       │
              └──────────┬──────────┘   └──────────┬──────────┘
                         │                      │
                         ▼                      ▼
              ┌─────────────────────┐   ┌─────────────────────┐
              │ Phase 8: Summary    │   │ Phase 5: TDD Red    │
              │ (完成总结)          │   │ (基于 Spec 写测试)  │
              └─────────────────────┘   └──────────┬──────────┘
                                                   │
                                                   ▼
                                    ┌─────────────────────┐
                                    │ Phase 6: TDD Green   │
                                    │ (实现代码)           │
                                    └──────────┬──────────┘
                                               │
                                               ▼
                                    ┌─────────────────────┐
                                    │ Phase 7: Refactor    │
                                    │ (重构+安全+审查)     │
                                    └──────────┬──────────┘
                                               │
                                               ▼
                                    ┌─────────────────────┐
                                    │ Phase 8: Summary    │
                                    │ (完成总结)           │
                                    └─────────────────────┘
```

### 简单流程 (SIMPLE)

**适用**: Bug 修复、小修改、单文件变更

| 阶段 | 名称 | Agent | 跳过 |
|:---|:---|:---|:---:|
| 0 | llmdoc 优先读取 | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 3 | Clarification | 主 Agent + AskUser | Phase 2, 4, 4.5 |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Security + Review | `code-simplifier` + `security-guidance` + `code-reviewer` | - |
| 8 | Summary | `investigator` | - |

**跳过的阶段**: Phase 2 (Exploration), Phase 4 (Architecture), Phase 4.5 (Spec)

### 完整流程 (COMPLEX)

**适用**: 新功能、模块开发、重构

| 阶段 | 名称 | Agent | 并行 |
|:---|:---|:---|:---:|
| 0 | llmdoc 优先读取 + ace-tool | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 2 | Exploration | `investigator` x3 | 是 |
| 3 | Clarification | 主 Agent + AskUser | - |
| 4 | Architecture | `architect` x2 + `investigator` | 是 |
| **4.5** | **Spec** | `architect` + `ace-tool` | **-** |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Security + Review | `code-simplifier` + `security-guidance` + `code-reviewer` | 是 |
| 8 | Summary | `investigator` | - |

## 硬约束

1. **Phase 0 强制**: 必须优先读取 llmdoc，不存在时自动生成
2. **智能分流**: Phase 1 后自动判断任务类型
3. **Phase 3 不可跳过**: 必须澄清所有疑问
4. **Phase 5 TDD 强制**: 必须先写测试
5. **并行优先**: 能并行的任务必须并行执行
6. **状态持久化**: 每个阶段后自动更新状态文件（无需确认）
7. **选项式编程**: 使用 AskUserQuestion 展示选项
8. **文档更新询问**: Phase 7 必须询问

## 循环状态管理

**自动化**: 每个阶段后自动更新 `.claude/gclm.{task_id}.local.md` frontmatter，无需用户确认：

```yaml
current_phase: <下一阶段编号>
phase_name: "<下一阶段名称>"
```

**状态更新的自动化原因**:
- 状态文件是内部元数据，不是代码
- 更新是确定性的（阶段完成 → 状态更新）
- 不影响代码质量或安全性

**仍需授权的场景**:
- Phase 4: Architecture 设计方案审批
- Phase 7: 文档更新询问

当所有 8 阶段完成，输出完成信号：
```
<promise>GCLM_WORKFLOW_COMPLETE</promise>
```

提前退出：在状态文件中设置 `active: false`。

---

## Phase 0: llmdoc 优先读取 + auggie 搜索

### 自动化流程

1. **检查 auggie 是否可用**
   - 运行 `auggie --help` 检查是否安装
   - 不可用 → 提示安装：`npm install -g @augmentcode/auggie@prerelease`

2. **检查 llmdoc/ 是否存在**
   - 存在 → 直接读取
   - 不存在 → **自动生成（不需要用户确认，直接执行）**

3. **自动生成 llmdoc（NON-NEGOTIABLE - 无需确认）**
   - 使用 `investigator` agent 扫描代码库
   - 生成 `llmdoc/index.md`
   - 生成 `llmdoc/overview/` 基础文档（project.md, tech-stack.md, structure.md）
   - **注意：这是初始化步骤，自动执行，不要询问用户**

4. **继续读取流程**
   - 读取 `llmdoc/index.md`
   - 读取 `llmdoc/overview/*.md` 全部
   - 根据任务读取 `llmdoc/architecture/*.md`

5. **auggie 搜索增强（可选）**
   - 当需要查找特定代码时使用 auggie MCP 的上下文搜索工具
   - 支持自然语言代码搜索
   - 自动从 IDE 获取项目路径

### auggie 工作原理

**重要**: auggie 是 MCP 服务器，提供高级代码上下文搜索：

```bash
# 安装 auggie
npm install -g @augmentcode/auggie@prerelease

# 配置环境变量（可选）
# AUGMENT_API_TOKEN: API 令牌
# AUGMENT_API_URL: API 端点

# MCP 配置已在 gclm-flow 的 .mcp.json 中配置
```

### 生成约束

- **最小化生成**: 只生成基础文档
- **增量完善**: 后续可在 Phase 7 补充
- **保持简洁**: 避免过度生成
- **直接执行**: llmdoc 不存在时自动生成，**不询问用户**

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

### Phase 4: Architecture (2 个并行方案 + 1 个测试策略)

**重要**: 必须等待 agents 完成并展示方案后，再询问用户选择

```bash
# 步骤 1: 并行启动 3 个 agents
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

---TASK---
id: p4_test_strategy
agent: gclm-investigator
workdir: .
---CONTENT---
Analyze testing strategy for this change.
EOF

# 步骤 2: 等待完成后，使用 TaskOutput 获取每个 agent 的输出
TaskOutput("p4_minimal", block=true)
TaskOutput("p4_pragmatic", block=true)
TaskOutput("p4_test_strategy", block=true)

# 步骤 3: 格式化展示方案给用户
# (将 3 个方案以清晰的格式展示)

# 步骤 4: 等待用户阅读后，使用 AskUserQuestion 询问选择
```

**关于 llmdoc**: Phase 4 不会自动生成/更新 llmdoc，文档更新在 Phase 7 询问用户后进行

## Agent 体系

| Agent | 职责 | 模型 |
|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku 4.5 |
| `architect` | 架构设计、方案权衡 | Opus 4.5 |
| `worker` | 执行明确定义的任务 | Sonnet 4.5 |
| `tdd-guide` | TDD 流程指导 | Sonnet 4.5 |
| `code-simplifier` | 代码简化重构 | Sonnet 4.5 |
| `security-guidance` | 安全审查 | Sonnet 4.5 |
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
- Tdd-guide output: <paste or "None">
- Code-simplifier output: <paste or "None">
- Security-guidance output: <paste or "None">
- Code-reviewer output: <paste or "None">
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

---

## auggie 快速参考

### 安装
```bash
# 全局安装 auggie
npm install -g @augmentcode/auggie@prerelease
```

### MCP 工具
Claude Code 可直接调用 auggie 提供的 MCP 工具进行：
- 自然语言代码搜索
- 代码上下文增强
- 语义代码理解

### 使用示例
```javascript
// Claude Code 自动调用，无需手动命令
// 搜索 "用户认证相关的代码"
// auggie 会自动理解意图并返回相关代码片段和上下文
```

### 配置
```bash
# 环境变量（可选）
export AUGMENT_API_TOKEN="your-token"
export AUGMENT_API_URL="https://acemcp.heroman.wtf/relay/"
```

### 项目支持
auggie 支持多种编程语言和文件类型，提供智能代码搜索和上下文理解。
