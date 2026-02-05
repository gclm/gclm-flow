# 工作流配置

## 概述

gclm-flow 的工作流通过 YAML 文件定义，位于 `workflows/` 目录。每个工作流定义了一组按依赖关系执行的 Agent 节点。

---

## 工作流类型

| 类型 | workflow_type | 适用场景 | 阶段数 |
|:---|:---|:---|:---:|
| 🔍 **ANALYZE** | `ANALYZE` | 代码分析、问题诊断、性能评估 | 5+1 |
| 📝 **DOCUMENT** | `DOCUMENT` | 文档编写、架构设计、需求分析 | 7+1 |
| 🔧 **CODE_SIMPLE** | `CODE_SIMPLE` | Bug 修复、小修改、单文件变更 | 6+1 |
| 🚀 **CODE_COMPLEX** | `CODE_COMPLEX` | 新功能、模块开发、跨文件变更 | 9+1 |

> **+1** = 可选的 `doc_update` 阶段（文档更新）

---

## YAML 结构

```yaml
name: workflow_name                    # 工作流唯一标识
display_name: "显示名称"                   # 人类可读名称
description: "工作流描述"                   # 详细说明
version: "1.0"                            # 版本号
author: "作者"                              # 作者
workflow_type: "CODE_SIMPLE"             # 工作流类型

nodes:                                    # 节点列表
  - ref: phase_id                          # 节点唯一标识
    display_name: "阶段名称"               # 显示名称
    agent: investigator                   # 使用的 Agent
    model: haiku                          # 使用的模型 (haiku/sonnet/opus)
    timeout: 60                            # 超时时间（秒）
    required: true                         # 是否必需
    depends_on:                            # 依赖节点
      - previous_phase
    parallel_group: ""                    # 并行组（可选）
    config:                                # 额外配置（可选）
      key: value

completion:                               # 完成配置（可选）
  signal: "<promise>GCLM_WORKFLOW_COMPLETE</promise>"
  final_status: completed

error_handling:                           # 错误处理（可选）
  max_retries: 1
  retry_on: [timeout, api_error]
  continue_on_non_required: true
```

---

## 内置工作流

### analyze.yaml

**用途**: 代码分析、问题诊断、性能评估、安全审计

**节点**:
1. `discovery` - 需求发现 (investigator, haiku)
2. `analysis` (并行组):
   - `code_analysis` - 代码分析 (investigator, sonnet)
   - `dependency_analysis` - 依赖分析 (investigator, sonnet) **[可选]**
   - `performance_analysis` - 性能分析 (investigator, sonnet) **[可选]**
3. `security_review` - 安全审查 (security-guidance, sonnet) **[可选]**
4. `report` - 分析报告 (investigator, sonnet)
5. `doc_update` - 文档更新 (llmdoc, sonnet) **[可选]**

### code_simple.yaml

**用途**: Bug 修复、小修改、单文件变更

**节点**:
1. `discovery` - 需求发现 (investigator, haiku)
2. `clarification` - 澄清确认 (investigator, haiku)
3. `tdd_red` - TDD Red (tdd-guide, sonnet)
4. `tdd_green` - TDD Green (worker, sonnet)
5. `review` (并行组):
   - `code_simplifier` - 代码简化
   - `security_guidance` - 安全审查
   - `code_reviewer` - 代码审查
6. `summary` - 完成总结 (investigator, haiku)
7. `doc_update` - 文档更新 (llmdoc, sonnet) **[可选]**

### code_complex.yaml

**用途**: 新功能、模块开发、跨文件变更

**节点**:
1. `discovery` - 需求发现 (investigator, haiku)
2. `exploration` - 探索研究 (investigator, haiku)
3. `clarification` - 澄清确认 (investigator, haiku)
4. `architecture` - 架构设计 (architect, opus)
5. `spec` - 规范文档 (spec-guide, opus)
6. `tdd_red` - TDD Red (tdd-guide, sonnet)
7. `tdd_green` - TDD Green (worker, sonnet)
8. `refactor_review` - 重构审查 (code-reviewer, sonnet)
9. `summary` - 完成总结 (investigator, haiku)
10. `doc_update` - 文档更新 (llmdoc, sonnet) **[可选]**

### document.yaml

**用途**: 文档编写、架构设计、需求分析

**节点**:
1. `discovery` - 需求发现 (investigator, haiku)
2. `exploration` - 探索研究 (investigator, haiku)
3. `clarification` - 澄清确认 (investigator, haiku)
4. `draft` - 起草文档 (architect, opus)
5. `refine` - 完善内容 (worker, sonnet)
6. `review` - 质量审查 (code-reviewer, sonnet)
7. `summary` - 完成总结 (investigator, haiku)
8. `doc_update` - 文档更新 (llmdoc, sonnet) **[可选]**

---

## 节点依赖

### 串行依赖

```yaml
nodes:
  - ref: phase_a
  - ref: phase_b
    depends_on: [phase_a]  # phase_b 在 phase_a 完成后执行
```

### 并行执行

```yaml
nodes:
  - ref: review_1
    parallel_group: review   # 与同组节点并行
  - ref: review_2
    parallel_group: review   # 与 review_1 并行
  - ref: review_3
    parallel_group: review   # 与 review_1, review_2 并行
```

### 混合依赖

```yaml
nodes:
  - ref: phase_1
  - ref: phase_2a
    depends_on: [phase_1]
    parallel_group: group_a
  - ref: phase_2b
    depends_on: [phase_1]
    parallel_group: group_a
  - ref: phase_3
    depends_on: [phase_2a, phase_2b]  # 等待 group_a 全部完成
```

---

## 工作流加载

### 内置工作流

`internal/db/workflow.go` 中的 `InitializeBuiltinWorkflows` 函数会自动加载内置工作流：

```go
builtinWorkflows := []struct {
    file   string
    name   string
    wtype  string
}{
    {"document.yaml", "document", "document"},
    {"code_simple.yaml", "code_simple", "code_simple"},
    {"code_complex.yaml", "code_complex", "code_complex"},
}
```

### 自定义工作流

通过 CLI 命令安装：

```bash
gclm-engine workflow install /path/to/custom_workflow.yaml
```

或直接复制到工作流目录：

```bash
cp custom_workflow.yaml ~/.gclm-flow/workflows/
```

---

## 工作流解析

`internal/pipeline/parser.go` 提供工作流解析功能：

### 主要方法

| 方法 | 功能 |
|:---|:---|
| `LoadPipeline(name)` | 加载指定工作流 |
| `LoadAllPipelines()` | 加载所有工作流 |
| `ValidatePipeline(pipeline)` | 验证工作流配置 |
| `GetPipelineByWorkflowType(type)` | 按类型获取工作流 |
| `CalculateExecutionOrder(pipeline)` | 计算执行顺序（含并行组） |

### 验证规则

1. **必需字段检查**: name, display_name, workflow_type, nodes
2. **依赖验证**: 依赖的节点必须存在
3. **循环依赖检查**: 使用 Kahn 算法检测 DAG

---

## 执行顺序

Go 引擎使用拓扑排序计算节点执行顺序：

```go
// 返回类型
type NodeExecutionOrder struct {
    Node     *PipelineNode
    Order    int    // 执行顺序
    Parallel int    // 并行组编号 (>0 表示并行)
}
```

**示例**:

```yaml
nodes:
  - ref: a        # Order: 0, Parallel: 0
  - ref: b        # Order: 1, Parallel: 1
    parallel_group: g1
  - ref: c        # Order: 1, Parallel: 1
    parallel_group: g1
  - ref: d        # Order: 2, Parallel: 0
    depends_on: [b, c]
```

---

## 自定义工作流

### 创建步骤

1. 复制示例模板：
   ```bash
   cp workflows/examples/custom_simple.yaml my_workflow.yaml
   ```

2. 编辑工作流配置：
   ```yaml
   name: my_workflow
   workflow_type: CODE_SIMPLE
   nodes:
     - ref: discovery
       agent: investigator
       model: haiku
       # ...
   ```

3. 安装工作流：
   ```bash
   gclm-engine workflow install my_workflow.yaml
   ```

4. 验证工作流：
   ```bash
   gclm-engine workflow validate my_workflow.yaml
   ```

### 工作流模板

`workflows/examples/` 目录提供了三个模板：

| 文件 | 类型 | 用途 |
|:---|:---|:---|
| `custom_simple.yaml` | CODE_SIMPLE | 最小化工作流模板 |
| `custom_document.yaml` | DOCUMENT | 文档编写模板 |
| `custom_complex.yaml` | CODE_COMPLEX | 完整 SpecDD 模板 |
