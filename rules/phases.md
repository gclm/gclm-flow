# gclm-flow 阶段规则

## 智能分流工作流

### 核心理念：SpecDD + TDD

**SpecDD** (Specification-Driven Development) 用于复杂模块开发，**TDD** (Test-Driven Development) 用于简单功能修复。

---

## 自动分类逻辑

### Phase 1 后自动判断任务类型

```python
# 智能分类伪代码
def classify_task(discovery_output, user_request):
    score = 0

    # 关键词分析
    simple_keywords = ["bug", "修复", "error", "fix", "问题", "调试"]
    complex_keywords = ["功能", "模块", "新", "开发", "重构", "系统", "设计"]

    for kw in simple_keywords:
        if kw in user_request.lower():
            score -= 2

    for kw in complex_keywords:
        if kw in user_request.lower():
            score += 2

    # 文件数量
    estimated_files = discovery_output.get("estimated_files", 1)
    if estimated_files <= 2:
        score -= 1
    elif estimated_files >= 5:
        score += 2

    # 风险评估
    if discovery_output.get("risk") == "high":
        score += 1

    # 分类
    if score <= -2:
        return "SIMPLE"   # 简单任务
    elif score >= 2:
        return "COMPLEX"  # 复杂任务
    else:
        return "MEDIUM"   # 需要用户确认
```

### 分类结果处理

| 分类 | 流程 | 适用场景 |
|:---|:---|:---|
| **SIMPLE** | 简单流程 | Bug 修复、单个函数修改、小幅重构 |
| **MEDIUM** | 询问用户 | 边界情况，让用户选择 |
| **COMPLEX** | 完整流程 | 新功能开发、模块重写、跨文件变更 |

---

## 简单流程 (SIMPLE)

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

---

## 完整流程 (COMPLEX)

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

---

## 阶段详细规则

### Phase 0: llmdoc 优先读取 + auggie (NON-NEGOTIABLE)

**目标**: 在任何代码操作前建立上下文理解

**步骤**:
1. **检查 auggie 是否可用**
   - 运行 `auggie --help` 检查是否安装
   - 不可用 → 提示安装：`npm install -g @augmentcode/auggie@prerelease`
2. 检查 `llmdoc/` 是否存在
3. **如果存在**:
   - 读取 `llmdoc/index.md` 获取导航
   - 读取 `llmdoc/overview/*.md` 全部文档
   - 根据任务读取 `llmdoc/architecture/*.md`
4. **如果不存在**:
   - 使用 `investigator` agent 扫描代码库
   - 自动生成 `llmdoc/index.md`
   - 自动生成 `llmdoc/overview/` 基础文档
   - 然后读取生成的文档
5. **auggie 搜索增强（可选）**
   - 需要查找特定代码时使用 auggie MCP 的上下文搜索

**输出**: 上下文摘要（关键文件、模块依赖、设计模式）

**强制**: 此阶段不可跳过
**自动化**: llmdoc 不存在时自动生成，无需用户确认

**auggie 使用**:
```bash
# 安装
npm install -g @augmentcode/auggie@prerelease

# MCP 配置自动生效
# Claude Code 可直接调用上下文搜索工具
```

---

### Phase 1: Discovery - 理解需求 + 任务分类

**Agent**: `investigator`

**输出**:
- Requirements (需求)
- Non-goals (非目标)
- Risks (风险)
- Acceptance Criteria (验收标准)
- **Task Classification** (任务分类: SIMPLE/COMPLEX/MEDIUM)
- Estimated Files (预估文件数)

**关键问题**:
- 用户可见行为是什么？
- 范围和边界在哪里？
- 有哪些约束和限制？
- 成功的标准是什么？
- 预估涉及多少文件？

**分类信号**:
```yaml
simple_signals:
  keywords: ["bug", "修复", "error", "fix", "问题", "调试"]
  file_count: "<= 2"
  risk: "low"

complex_signals:
  keywords: ["功能", "模块", "新", "开发", "重构", "系统", "设计"]
  file_count: ">= 5"
  risk: "any"
```

---

### Phase 2: Exploration - 探索代码库

**并行执行 3 个 `investigator`**

| 任务 | 描述 | 输出 |
|:---|:---|:---|
| 相似功能 | 查找 1-3 个相似功能 | 关键文件、调用流程、扩展点 |
| 架构映射 | 映射相关子系统 | 模块图 + 5-10 个关键文件 |
| 代码规范 | 识别测试模式、规范 | 测试命令 + 文件位置 |

**并行执行**: 必须在单个响应中使用多个 Task 调用

**auggie 集成**: 使用语义搜索加速代码探索

---

### Phase 3: Clarification - 澄清疑问 (强制)

**不可跳过的阶段**

1. 汇总 Phase 1 和 Phase 2 输出
2. **如果分类为 MEDIUM，询问用户选择流程**
3. 生成优先级排序的问题列表
4. 使用 `AskUserQuestion` 逐一确认

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

**工作流程** (重要：顺序执行，不可跳过)：

1. **等待 agents 完成** - 等待 3 个并行任务全部完成
2. **收集并展示方案** - 使用 `TaskOutput` 获取每个 agent 的完整输出
3. **格式化展示** - 将 3 个方案以清晰的格式展示给用户
4. **等待用户阅读** - 给用户时间阅读和比较方案
5. **使用 `AskUserQuestion` 选择** - 用户阅读后再询问选择

**显式审批门**: "Approve starting implementation?"

**输出应包含**:
- 文件清单（创建/修改）
- 组件设计
- 数据流
- 构建序列

**关于 llmdoc**:
- Phase 4 **不会**自动生成/更新 llmdoc
- llmdoc 更新在 **Phase 7** 询问用户确认后才会进行

---

### Phase 4.5: Spec - 编写规范文档 (SpecDD)

**目标**: 为复杂模块编写详细的规范文档

**Agent**: `architect` + `auggie`

**输入**: Phase 4 的设计方案

**输出**: `.claude/specs/{feature-name}.md`

**Spec 文档结构**:
```markdown
# {功能名称} 规范文档

## 1. 概述
### 1.1 目标
### 1.2 范围
### 1.3 非目标

## 2. 功能需求
### 2.1 用户故事
### 2.2 验收标准
### 2.3 边界条件

## 3. API 设计
### 3.1 公开接口
### 3.2 数据结构
### 3.3 错误处理

## 4. 技术设计
### 4.1 组件架构
### 4.2 数据流
### 4.3 依赖关系

## 5. 测试策略
### 5.1 单元测试覆盖
### 5.2 集成测试场景
### 5.3 边界测试

## 6. 非功能需求
### 6.1 性能要求
### 6.2 安全要求
### 6.3 可维护性
```

**auggie 集成**: 使用语义搜索查找相关代码示例

---

### Phase 5: TDD Red - 编写测试

**Agent**: `tdd-guide`

**完整流程**: 基于 Phase 4.5 的 Spec 编写测试
**简单流程**: 基于 Phase 1 的理解直接编写测试

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

### Phase 7: Refactor + Security + Review - 重构、安全与审查

**并行执行**:

| Agent | 任务 |
|:---|:---|
| code-simplifier | 代码简化 - 清晰度、一致性、可维护性 |
| security-guidance | 安全审查 - 漏洞检测、安全最佳实践 |
| code-reviewer | 代码审查 - 正确性 + 简洁性 |

**重构原则**:
- 保持测试绿色
- 消除重复
- 改进命名
- 优化性能
- 修复安全隐患

**文档更新询问**:
```
AskUserQuestion: "是否更新项目文档 (llmdoc 和 Spec)？"
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
workflow_type: "simple"  # simple | complex
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

**仍需授权的场景**:
- Phase 3: 流程类型选择 (MEDIUM 时)
- Phase 4: Architecture 设计方案审批
- Phase 7: 文档更新询问

---

## 并行执行模式

### 必须并行的阶段 (完整流程)

- **Phase 2**: 3 个 investigator
- **Phase 4**: 2 个 architect + 1 个 investigator
- **Phase 7**: code-simplifier + security-guidance + code-reviewer

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
