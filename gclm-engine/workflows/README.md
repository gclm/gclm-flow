# 自定义工作流示例

本目录包含 gclm-engine 的自定义工作流示例，可作为创建自己工作流的模板。

## 工作流类型

gclm-engine 默认支持三种工作流类型：

| 类型 | workflow_type | 适用场景 |
|:---|:---|:---|
| 📝 **DOCUMENT** | `DOCUMENT` | 文档编写、方案设计、需求分析 |
| 🔧 **CODE_SIMPLE** | `CODE_SIMPLE` | Bug 修复、小修改、单文件变更 |
| 🚀 **CODE_COMPLEX** | `CODE_COMPLEX` | 新功能、模块开发、跨文件变更 |

## 示例文件

### code_simple.yaml
最简化的工作流，包含：
- 需求发现
- 澄清确认
- 实现
- 总结

适合快速任务和 Bug 修复。

### document.yaml
文档专用工作流，包含：
- 需求发现
- 探索研究
- 澄清确认
- 起草文档
- 完善内容
- 质量审查
- 完成总结

适合编写技术文档、API 文档、设计方案。

### code_complex.yaml
完整 SpecDD 工作流，包含：
- 需求发现
- 探索研究
- 澄清确认
- 架构设计
- 规范文档 (Spec)
- TDD Red (测试)
- TDD Green (实现)
- 重构审查
- 完成总结

适合复杂功能开发和跨模块变更。

## 如何使用

### 1. 复制示例作为起点

```bash
cp code_simple.yaml my_workflow.yaml
```

### 2. 编辑工作流

修改 `my_workflow.yaml` 中的以下字段：
- `name`: 工作流唯一标识符
- `display_name`: 显示名称
- `description`: 描述
- `workflow_type`: 工作流类型
- `nodes`: 添加/修改阶段节点

### 3. 安装工作流

```bash
# 方法 1: 使用 gclm-engine 安装
~/.gclm-flow/gclm-engine workflow install my_workflow.yaml

```

### 4. 使用工作流

```bash
# 列出所有工作流
~/.gclm-flow/gclm-engine workflow list

# 使用自定义工作流创建任务
~/.gclm-flow/gclm-engine workflow start "你的任务描述" --workflow my_workflow
```

## 节点配置

每个节点包含以下配置项：

```yaml
- ref: phase_id           # 节点唯一标识符
  display_name: "显示名称"  # 显示名称
  agent: investigator     # 使用的 Agent
  model: haiku           # 使用的模型 (haiku/sonnet/opus)
  timeout: 60            # 超时时间（秒）
  required: true         # 是否必需
  depends_on:            # 依赖节点列表
    - previous_phase
  config:                # 额外配置（可选）
    key: value
```

## 可用的 Agents

| Agent | 模型 | 用途 |
|:---|:---|:---:|
| `investigator` | Haiku | 探索、分析、总结 |
| `architect` | Opus | 架构设计、方案权衡 |
| `spec-guide` | Opus | SpecDD 规范文档编写 |
| `tdd-guide` | Sonnet | TDD 流程指导 |
| `worker` | Sonnet | 执行明确定义的任务 |
| `code-reviewer` | Sonnet | 代码审查 |

## 节点依赖

使用 `depends_on` 定义节点间的依赖关系：

```yaml
- ref: phase_a
  # ... 其他配置

- ref: phase_b
  depends_on:
    - phase_a  # phase_b 会在 phase_a 完成后执行
```

## 并行执行

使用 `parallel_group` 实现并行执行：

```yaml
- ref: review_1
  parallel_group: review  # 与同组节点并行执行

- ref: review_2
  parallel_group: review  # 与 review_1 并行
```

## 最佳实践

1. **保持简洁**: 从简单的工作流开始，逐步扩展
2. **合理超时**: 根据任务复杂度设置合适的超时时间
3. **明确依赖**: 使用 `depends_on` 确保阶段按正确顺序执行
4. **使用必需标志**: 关键阶段设置 `required: true`
5. **文档清晰**: 为工作流和节点提供清晰的描述

## 更多信息

- [gclm-engine 文档](../README.md)
