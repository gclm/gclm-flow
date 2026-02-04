# 配置参考

## 概览

gclm-flow 的配置分布在多个文件中：

| 文件 | 位置 | 用途 |
|:---|:---|:---|
| `CLAUDE.md` | `~/.claude/` | 核心配置和规则 |
| `settings.json` | `~/.claude/` | MCP 服务器配置 |
| `settings.local.json` | `~/.claude/` | 本地覆盖配置 |
| `.claude/gclm.*.local.md` | 项目目录 | 工作流状态 |

---

## CLAUDE.md

### 核心配置

`CLAUDE.md` 包含 gclm-flow 的核心配置，通常通过复制示例文件：

```bash
cp CLAUDE.example.md ~/.claude/CLAUDE.md
```

### 主要部分

1. **gclm-flow 核心配置** - 工作流哲学和 Agent 体系
2. **Agent 编排规则** - Agent 调用时机和并行策略
3. **llmdoc 文档规则** - 文档优先策略
4. **阶段规则** - 工作流阶段详细说明
5. **TDD 规范** - 测试驱动开发规范
6. **SpecDD 规范** - 规范驱动开发规范

### 修改配置

编辑 `~/.claude/CLAUDE.md` 来自定义工作流行为。

---

## settings.json

### MCP 服务器配置

`settings.json` 定义 MCP (Model Context Protocol) 服务器：

```json
{
  "mcpServers": {
    "auggie": {
      "command": "npx",
      "args": ["@augmentcode/auggie@prerelease"]
    }
  }
}
```

### auggie 配置

**描述**: 语义代码搜索服务器

**命令**: `npx @augmentcode/auggie@prerelease`

**环境变量** (可选):
```bash
export AUGMENT_API_TOKEN="your-token"
export AUGMENT_API_URL="https://acemcp.heroman.wtf/relay/"
```

---

## settings.local.json

### 本地覆盖配置

`settings.local.json` 用于本地配置覆盖，通常包含：

1. **权限配置** - 自动授权的 Bash 命令
2. **Hooks 配置** - 生命周期钩子
3. **本地 MCP 服务器** - 项目特定的 MCP 服务

### 示例

```json
{
  "permissions": {
    "allow": [
      "Bash(\"mkdir -p .claude*\")",
      "Bash(\"ls -la .claude/\")",
      "Bash(\"cat .claude/gclm.*.local.md\")"
    ]
  },
  "hooks": {
    "PermissionRequest": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/hooks/notify.sh '需要权限确认'"
          }
        ]
      }
    ],
    "Notification": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.claude/hooks/notify.sh '有待处理的操作'"
          }
        ]
      }
    ]
  }
}
```

---

## 工作流状态文件

### 位置

`.claude/gclm.{task_id}.local.md`

### 结构

```yaml
---
active: true
current_phase: 0
phase_name: "llmdoc Reading"
max_phases: 9
workflow_type: "CODE_COMPLEX"
task_description: "任务描述"
completion_promise: "<promise>GCLM_WORKFLOW_COMPLETE</promise>"
---
```

### 字段说明

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `active` | boolean | 工作流是否活跃 |
| `current_phase` | number | 当前阶段 (0-9) |
| `phase_name` | string | 阶段名称 |
| `max_phases` | number | 最大阶段数 |
| `workflow_type` | string | 工作流类型 (DOCUMENT/CODE_SIMPLE/CODE_COMPLEX) |
| `task_description` | string | 任务描述 |
| `completion_promise` | string | 完成信号 |

### 手动管理

**查看状态**:
```bash
cat .claude/gclm.*.local.md
```

**强制退出**:
```bash
sed -i.bak 's/^active: true/active: false/' .claude/gclm.*.local.md
```

---

## 环境变量

### CLAUDE_DIR

**默认**: `~/.claude`

**用途**: Claude Code 配置目录

**示例**:
```bash
export CLAUDE_DIR="$HOME/.claude"
```

### PROJECT_DIR

**默认**: 自动检测

**用途**: gclm-flow 项目目录

**示例**:
```bash
export PROJECT_DIR="/path/to/gclm-flow"
```

### AUGMENT_API_TOKEN

**用途**: auggie API 令牌

**示例**:
```bash
export AUGMENT_API_TOKEN="your-token-here"
```

### AUGMENT_API_URL

**默认**: `https://acemcp.heroman.wtf/relay/`

**用途**: auggie API 服务器

**示例**:
```bash
export AUGMENT_API_URL="https://your-auggie-server.com/relay/"
```

---

## Hooks 配置

### Hook 类型

| 类型 | 触发时机 | 配置键 |
|:---|:---|:---|
| `PermissionRequest` | 需要权限确认时 | `hooks.PermissionRequest` |
| `Notification` | 有待处理的操作时 | `hooks.Notification` |
| `Stop` | 用户尝试退出时 | `hooks.Stop` |

### notify.sh

**位置**: `~/.claude/hooks/notify.sh`

**功能**: 通知用户有权限请求或待处理操作

**权限**: 必须可执行 (`chmod +x`)

### stop-gclm-loop.sh

**位置**: `~/.claude/hooks/stop-gclm-loop.sh`

**功能**: 阻止活跃工作流的中途退出

**状态**: 未在默认配置中注册（可选）

---

## 模型选择

### Agent 模型映射

| Agent | 模型 | 原因 |
|:---|:---|:---|
| `investigator` | Haiku 4.5 | 速度快，成本低 |
| `architect` | Opus 4.5 | 深度思考，高质量 |
| `spec-guide` | Opus 4.5 | 复杂规范需要深度思考 |
| `tdd-guide` | Sonnet 4.5 | 平衡速度和质量 |
| `worker` | Sonnet 4.5 | 标准实现 |
| `code-reviewer` | Sonnet 4.5 | 代码审查 |

### 自定义模型

如需使用不同模型，编辑相应的 Agent 定义文件：

```markdown
---
model: opus  # 或 sonnet, haiku
---
```

---

## 配置验证

### 检查安装

```bash
# 检查 agents
ls ~/.claude/agents/

# 检查 commands
ls ~/.claude/commands/

# 检查 skills
ls ~/.claude/skills/

# 检查 hooks
ls ~/.claude/hooks/

# 检查权限
ls -la ~/.claude/hooks/
```

### 检查配置

```bash
# 检查 CLAUDE.md
cat ~/.claude/CLAUDE.md

# 检查 settings.json
cat ~/.claude/settings.json

# 检查 MCP 服务器
cat ~/.claude/settings.local.json | grep mcpServers
```

### 检查 auggie

```bash
# 检查安装
which auggie

# 检查版本
auggie --version

# 测试搜索
auggie "用户认证"
```

---

## 故障排查

### 问题: 配置未生效

**解决**:
1. 检查文件路径是否正确
2. 检查文件权限
3. 重启 Claude Code

---

### 问题: MCP 服务器无法启动

**解决**:
1. 检查 `settings.json` 语法
2. 检查命令是否可执行
3. 查看错误日志

---

### 问题: Hooks 不工作

**解决**:
1. 检查脚本权限 (`chmod +x`)
2. 检查 shebang (`#!/bin/bash`)
3. 手动测试脚本
