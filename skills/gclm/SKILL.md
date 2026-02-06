---
name: gclm
description: "智能分流工作流引擎 - LLM 语义匹配选择工作流，支持动态扩展"
allowed-tools: [
  "Bash(gclm-engine *)",
  "Read(*)", "Write(*)", "Edit(*)",
  "Glob(*)", "Grep(*)", "Task(*)"
]
version: "5.0"
engine: "gclm-engine Go Engine"
---

## 核心设计

**动态工作流引擎** - 通过 Go 引擎 + YAML 配置实现可扩展的工作流系统：

```
用户请求 → LLM 语义分析 → 查询工作流列表 → 匹配最佳工作流 → 用户确认 → 执行
```

### 关键特性

- **动态扩展**: 在 `workflows/` 添加 YAML 文件即可新增工作流
- **LLM 语义匹配**: 由 LLM 根据任务描述和工作流描述进行语义匹配
- **用户可见**: 向用户展示选择的工作流和理由
- **用户可控**: 用户可以手动调整选择
- **状态持久化**: SQLite 数据库管理任务和阶段状态
- **JSON 输出**: 所有命令支持 `--json` 标志输出 JSON 格式

---

## 工作流类型

### 标准类型

| 类型 | display_name | 描述 |
|:---|:---|:---|
| `analyze` | 代码分析 | 代码分析、问题诊断、性能评估、架构分析 |
| `review` | 代码审查 | 代码审查、安全审计、质量检查 |
| `feat` | 新功能 | 新功能开发、模块开发、功能实现 |
| `fix` | Bug 修复 | Bug 修复、错误处理、问题解决 |
| `docs` | 文档 | 文档编写、方案设计、需求分析、API 文档 |
| `refactor` | 重构 | 代码重构、架构调整、优化改进 |
| `test` | 测试 | 测试编写、测试优化、覆盖率提升 |
| `chore` | 构建/工具 | 构建配置、工具升级、依赖更新 |
| `style` | 代码格式 | 代码格式调整、样式修改（不影响功能） |
| `perf` | 性能优化 | 性能优化、响应时间优化、资源优化 |
| `ci` | CI 配置 | CI/CD 配置、自动化脚本 |
| `deploy` | 部署 | 部署配置、发布流程、环境配置 |

### 类型定义位置

配置文件: `~/.gclm-flow/gclm_engine_config.yaml` (用户可扩展)

---

## 执行流程

### 步骤 1: 获取工作流列表

```bash
# 获取所有可用工作流（JSON 格式）
~/.gclm-flow/gclm-engine workflow list --json
```

**返回示例**:
```json
[
  {
    "name": "analyze",
    "displayName": "代码分析工作流",
    "description": "用于代码分析、问题诊断、性能评估、安全审计等纯分析任务",
    "workflowType": "analyze",
    "version": "1.0"
  },
  {
    "name": "docs",
    "displayName": "文档编写工作流",
    "description": "用于编写技术文档、设计方案、需求分析等文档类任务",
    "workflowType": "docs",
    "version": "1.0"
  },
  {
    "name": "feat",
    "displayName": "复杂功能开发工作流",
    "description": "用于新功能开发、模块开发、跨文件重构等复杂任务",
    "workflowType": "feat",
    "version": "1.0"
  },
  {
    "name": "fix",
    "displayName": "Bug修复工作流",
    "description": "Bug修复、小修改、单文件变更的标准流程",
    "workflowType": "fix",
    "version": "0.1.0-poc"
  }
]
```

### 步骤 2: LLM 语义匹配

根据用户任务描述和工作流描述进行语义匹配：

**分析逻辑**:
1. 提取任务关键词（如 "分析"、"bug"、"功能"、"文档"）
2. 匹配工作流类型（analyze、fix、feat、docs 等）
3. 选择最匹配的工作流

**示例匹配**:
| 用户输入 | 关键词 | 匹配工作流 | workflow_type |
|:---|:---|:---|:---|
| "分析用户认证模块的安全性" | 分析、安全性 | analyze | `analyze` |
| "修复登录按钮样式问题" | 修复、问题 | fix | `fix` |
| "添加用户认证功能" | 添加、功能 | feat | `feat` |
| "编写 API 设计文档" | 编写、文档 | docs | `docs` |
| "重构数据访问层" | 重构 | feat | `feat` |

### 步骤 3: 向用户展示选择

```
📋 工作流选择分析

根据您的任务: "分析用户认证模块的安全性"

我为您选择: 🔍 analyze (代码分析工作流)

选择理由:
- 任务关键词: "分析"、"安全性"
- 工作流描述: "用于代码分析、问题诊断、性能评估、安全审计等纯分析任务"
- 类型匹配: analyze
```

### 步骤 4: 用户确认

使用 `AskUserQuestion` 让用户选择：

```
是否使用 analyze 工作流？

选项:
- ✅ 使用 analyze (推荐)
- 🔄 手动选择其他工作流
- ❌ 取消任务
```

### 步骤 5: 启动工作流

```bash
# 创建任务（使用 --workflow 指定工作流名称）
~/.gclm-flow/gclm-engine task create "<任务描述>" --workflow <name> --json
```

**返回示例**:
```json
{
  "task_id": "task-xxx",
  "workflow_type": "analyze",
  "workflow": "analyze",
  "status": "created",
  "current_phase": 0,
  "total_phases": 7,
  "message": "Task created successfully"
}
```

---

## 阶段循环

### 1. 获取当前阶段

```bash
~/.gclm-flow/gclm-engine task current <task-id> --json
```

### 2. 执行阶段

根据 `current_phase` 的 `agent` 和 `model` 调用相应 Agent

### 3. 完成阶段

```bash
~/.gclm-flow/gclm-engine task complete <task-id> <phase-id> --output "<阶段输出>" --json
```

### 4. 标记失败（可选）

```bash
~/.gclm-flow/gclm-engine task fail <task-id> <phase-id> --error "<错误信息>" --json
```

### 5. 重复步骤 1-4

直到所有阶段完成

---

## 状态查询

```bash
# 查看完整执行计划
~/.gclm-flow/gclm-engine task plan <task-id> --json

# 查看事件日志
~/.gclm-flow/gclm-engine task events <task-id> --json

# 列出所有任务
~/.gclm-flow/gclm-engine task list --json
```

---

## 命令参考

### workflow 命令

| 命令 | 说明 |
|:---|:---|
| `workflow list` | 列出所有工作流 |
| `workflow info <name>` | 显示工作流详情 |
| `workflow validate <file>` | 验证 YAML 配置 |
| `workflow install <file>` | 安装工作流 |
| `workflow sync [file]` | 同步 YAML 到数据库 |

### task 命令

| 命令 | 说明 |
|:---|:---|
| `task create <prompt> --workflow <name>` | 创建任务（使用指定工作流） |
| `task get <task-id>` | 获取任务详情 |
| `task list` | 列出所有任务 |
| `task current <task-id>` | 获取当前待执行阶段 |
| `task plan <task-id>` | 获取执行计划 |
| `task complete <task-id> <phase-id> --output <text>` | 完成阶段 |
| `task fail <task-id> <phase-id> --error <msg>` | 标记阶段失败 |
| `task phases <task-id>` | 显示任务阶段 |
| `task events <task-id>` | 显示任务事件 |

### 全局标志

| 标志 | 说明 |
|:---|:---|
| `--json, -j` | 输出 JSON 格式（便于脚本解析） |

---

## 工作流定义

### YAML 结构

```yaml
name: my_workflow                    # 工作流唯一标识（文件名）
workflow_type: "feat"               # 必需，使用标准类型
display_name: "我的工作流"           # 人类可读名称
description: "工作流描述"           # LLM 匹配时的重要依据
version: "1.0"                      # 版本号

nodes:
  - ref: discovery                  # 节点唯一标识
    display_name: "需求发现"
    agent: investigator            # Agent 名称
    model: haiku                   # 模型 (haiku/sonnet/opus)
    timeout: 60                     # 超时（秒）
    required: true                  # 是否必需
    depends_on:                     # 依赖节点
      - previous_phase

  - ref: clarification
    display_name: "澄清确认"
    agent: investigator
    model: haiku
    depends_on: [discovery]
```

### 依赖和并行

**串行依赖**:
```yaml
- ref: phase_b
  depends_on: [phase_a]  # phase_b 等待 phase_a 完成
```

**并行执行**:
```yaml
- ref: review_1
  parallel_group: review   # 与同组节点并行
- ref: review_2
  parallel_group: review
```

---

## 添加新工作流

### 步骤

1. **创建 YAML 文件**:
   ```bash
   # workflows/my_custom_workflow.yaml
   name: my_custom_workflow
   workflow_type: "feat"      # 使用标准类型
   display_name: "我的自定义工作流"
   description: "用于特定场景..."
   nodes:
     # ...
   ```

2. **同步到数据库**:
   ```bash
   ~/.gclm-flow/gclm-engine workflow sync workflows/my_custom_workflow.yaml
   ```

3. **验证工作流**:
   ```bash
   ~/.gclm-flow/gclm-engine workflow validate workflows/my_custom_workflow.yaml
   ```

---

## 硬约束

1. **workflow_type 必需**: 所有工作流必须声明合法的 workflow_type
2. **名称一致性**: 文件名、name、workflow_type 三者保持一致
3. **用户确认**: 必须向用户展示选择的工作流并确认
4. **状态持久化**: 每个阶段后调用引擎更新状态
5. **--workflow 必需**: 创建任务时必须指定工作流名称（analyze, docs, feat, fix）

---

## 代码搜索

### 分层回退策略

| 方法 | 优势 | 劣势 | 状态 |
|:---|:---|:---|:---:|
| **auggie** | 语义搜索、自然语言查询 | 需要外部服务 | 推荐 |
| **llmdoc** | 结构化文档、本地 | 覆盖范围有限 | 默认 |
| **Grep** | 完整代码搜索 | 速度较慢 | 备选 |

### 安装 auggie（可选）

```bash
npm install -g @augmentcode/auggie@prerelease
```
