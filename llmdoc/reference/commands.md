# 命令参考

## gclm-engine CLI 命令

### workflow - 工作流管理

```bash
# 列出所有工作流
gclm-engine workflow list

# 验证工作流配置
gclm-engine workflow validate <workflow.yaml>

# 安装自定义工作流
gclm-engine workflow install <workflow.yaml>

# 查看工作流信息
gclm-engine workflow info <workflow-name>

# 启动工作流（自动检测类型）
gclm-engine workflow start "<任务描述>"

# 启动指定类型的工作流
gclm-engine workflow start "<任务描述>" --workflow document
```

### task - 任务管理

```bash
# 查看当前阶段
gclm-engine task current <task-id>

# 查看所有阶段
gclm-engine task phases <task-id>

# 完成阶段
gclm-engine task complete <task-id> <phase-id> --output "输出结果"

# 失败阶段
gclm-engine task fail <task-id> <phase-id> --error "错误信息"

# 查看任务详情
gclm-engine task show <task-id>

# 列出任务
gclm-engine task list [--status completed]
```

### 其他命令

```bash
# 版本信息
gclm-engine version

# 帮助信息
gclm-engine help
gclm-engine help workflow
gclm-engine help task
```

---

## JSON 输出

所有命令支持 `--json` 标志输出 JSON 格式：

```bash
gclm-engine workflow list --json
gclm-engine task current <task-id> --json
```

**输出格式**:

```json
{
  "id": "task-uuid",
  "status": "running",
  "current_phase": 2,
  "total_phases": 6,
  "next_phase": {
    "id": "phase-uuid",
    "name": "clarification",
    "display_name": "Clarification / 澄清确认",
    "agent": "investigator",
    "model": "haiku",
    "timeout": 60
  }
}
```

---

## Claude Code Skills 命令

### /gclm - 智能分流工作流

```bash
/gclm <任务描述>
```

**功能**: 智能分流工作流，自动判断任务类型并选择最优开发流程

**工作流类型**:

| 类型 | 检测关键词 | 阶段数 |
|:---|:---|:---:|
| 🔍 ANALYZE | 分析、诊断、审计、评估、检查 | 5 |
| 📝 DOCUMENT | 文档、方案、设计、需求 | 7 |
| 🔧 CODE_SIMPLE | bug、修复、error | 6 |
| 🚀 CODE_COMPLEX | 功能、模块、开发 | 9 |

**示例**:
```
/gclm 分析用户认证模块的安全性
/gclm 添加用户认证功能
/gclm 修复登录按钮样式
/gclm 编写 API 设计文档
/gclm 重构数据访问层
```

### /investigate - 代码库调查

```bash
/investigate <问题>
```

**功能**: 快速代码库调查，使用 investigator agent 分析项目

**示例**:
```
/investigate 项目中如何处理用户认证？
/investigate 错误处理机制在哪里？
/investigate 数据库连接是怎么建立的？
```

### /tdd - 测试驱动开发

```bash
/tdd <功能>
```

**功能**: 测试驱动开发，遵循 Red-Green-Refactor 循环

**TDD 循环**:
```
Red (写测试) → Green (写实现) → Refactor (重构)
```

**绝对规则**:
1. 绝不一次性生成代码和测试
2. 先写测试，后写实现
3. 测试必须先失败
4. 覆盖率 > 80%

### /spec - 规范驱动开发

```bash
/spec <功能>
```

**功能**: 规范驱动开发，先写详细规范文档，再编写测试和实现

**适用场景**:
- 新功能开发
- 跨模块变更 (3+ 文件)
- API 设计
- 数据结构设计

### /llmdoc - 文档生成/更新

```bash
/llmdoc
```

**功能**: 自动生成或更新项目 llmdoc 文档

**行为**:
1. 检查 `llmdoc/` 是否存在
2. 存在 → 扫描代码库并更新文档
3. 不存在 → 生成基础文档

---

## 命令对比

| 命令 | 复杂度 | 适用场景 | Agent 使用 |
|:---|:---|:---|:---|
| `/gclm` | 自动 | 所有场景 | 全部 |
| `/investigate` | 低 | 代码理解 | investigator |
| `/tdd` | 中 | 功能实现 | tdd-guide + worker |
| `/spec` | 高 | 架构设计 | architect + spec-guide + tdd-guide + worker |
| `/llmdoc` | 低 | 文档更新 | investigator |

---

## 环境变量

| 变量 | 默认值 | 说明 |
|:---|:---|:---|
| `GCLM_ENGINE_WORKFLOWS_DIR` | `~/.gclm-flow/workflows` | 工作流目录 |
| `GCLM_ENGINE_DB_PATH` | `~/.gclm-flow/gclm-engine.db` | 数据库路径 |
| `GCLM_VERSION` | `latest` | 安装时使用的版本 |

---

## 退出码

| 退出码 | 含义 |
|:---|:---|
| 0 | 成功 |
| 1 | 一般错误 |
| 2 | 参数验证失败 |
| 3 | 工作流未找到 |
| 4 | 数据库错误 |
| 5 | 循环依赖检测 |
