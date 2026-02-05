# 命令参考

## gclm-engine CLI 命令

### init - 初始化配置

```bash
# 初始化配置（导出内置工作流和配置）
gclm-engine init

# 强制覆盖现有文件
gclm-engine init --force

# 静默初始化（无输出）
gclm-engine init --silent
```

### workflow - 工作流管理

```bash
# 列出所有工作流（从数据库）
gclm-engine workflow list
gclm-engine workflow list --json

# 验证工作流配置
gclm-engine workflow validate <workflow.yaml>

# 安装自定义工作流（复制到 workflows/）
gclm-engine workflow install <workflow.yaml> [--name <custom-name>]

# 卸载自定义工作流
gclm-engine workflow uninstall <workflow-name>

# 查看工作流信息
gclm-engine workflow info <workflow-name>

# 导出工作流到 YAML
gclm-engine workflow export <workflow-name> [output-file]

# 同步工作流 YAML 到数据库（草稿 → 正式）
gclm-engine workflow sync                           # 同步所有
gclm-engine workflow sync workflows/feat.yaml      # 同步单个
gclm-engine workflow sync --force                  # 强制同步

# 启动工作流
gclm-engine workflow start "<任务描述>" --workflow <name>
```

### task - 任务管理

```bash
# 创建任务（已废弃，使用 workflow start）
gclm-engine task create "<提示>" --workflow-type CODE_SIMPLE

# 查看任务详情
gclm-engine task get <task-id>

# 列出任务
gclm-engine task list [--status completed] [--limit 20]

# 查看当前阶段（下一步要执行的）
gclm-engine task current <task-id>
gclm-engine workflow next <task-id>  # 别名

# 查看执行计划（所有阶段）
gclm-engine task plan <task-id>

# 查看所有阶段
gclm-engine task phases <task-id>

# 查看事件日志
gclm-engine task events <task-id> [--limit 50]

# 完成阶段
gclm-engine task complete <task-id> <phase-id> --output "输出结果"

# 失败阶段
gclm-engine task fail <task-id> <phase-id> --error "错误信息"

# 更新阶段状态
gclm-engine task update <task-id> <phase-id> completed --output "..."
gclm-engine task update <task-id> <phase-id> failed --error "..."

# 导出状态文件（兼容旧版 skills）
gclm-engine task export <task-id> <output-file>

# 任务控制
gclm-engine task pause <task-id>
gclm-engine task resume <task-id>
gclm-engine task cancel <task-id>
```

### pipeline - 流水线管理（保留兼容）

```bash
# 列出流水线（实际列出 workflows）
gclm-engine pipeline list

# 查看流水线详情
gclm-engine pipeline get <name>

# 推荐流水线（已废弃，使用 workflow list --json）
gclm-engine pipeline recommend "<提示>"
```

**注意**: `pipeline` 命令保留向后兼容，内部已重命名为 `workflow`。

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
  "task_id": "task-uuid",
  "workflow": "document",
  "workflow_type": "DOCUMENT",
  "total_phases": 7,
  "current_phase": {
    "phase_id": "phase-uuid",
    "phase_name": "discovery",
    "display_name": "Discovery / 需求发现",
    "agent": "investigator",
    "model": "haiku",
    "sequence": 1,
    "required": true,
    "timeout": 60
  }
}
```

**workflow list 输出**:

```json
[
  {
    "name": "document",
    "display_name": "DOCUMENT 工作流",
    "description": "文档编写、架构设计、需求分析",
    "workflow_type": "DOCUMENT",
    "version": "1.0.0",
    "is_builtin": true
  }
]
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
| 📝 DOCUMENT | 文档、方案、设计、需求 | 7 |
| 🔧 CODE_SIMPLE | bug、修复、error | 6 |
| 🚀 CODE_COMPLEX | 功能、模块、开发 | 9 |
| 🔍 ANALYZE | 分析、诊断、审计、评估 | 5 |

**流程**:
1. 调用 `workflow list --json` 获取所有工作流
2. LLM 根据提示语义选择最匹配的工作流
3. 调用 `workflow start "<提示>" --workflow <name>`

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

**推荐**:
- 新用户: 使用 `/gclm` (自动选择工作流)
- 高级用户: 直接调用 `workflow start --workflow <name>`
- 文档更新: 使用 `/llmdoc` 自动生成/更新文档

---

## 环境变量

| 变量 | 默认值 | 说明 |
|:---|:---|:---|
| `GCLM_ENGINE_CONFIG_DIR` | `~/.gclm-flow` | 配置目录 |
| `GCLM_ENGINE_DB_PATH` | `~/.gclm-flow/gclm-engine.db` | 数据库路径 |
| `GCLM_ENGINE_WORKFLOWS_DIR` | `~/.gclm-flow/workflows` | 工作流草稿目录 |
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
