# gclm-flow - 智能分流工作流插件

> SpecDD + TDD + llmdoc 优先 + auggie + 多 Agent 并行

<div align="center">

**智能分流工作流 - 自动判断任务类型，选择最优开发流程**

</div>

---

## 核心哲学

1. **SpecDD + TDD**: 复杂模块走 SpecDD，简单修改走 TDD
2. **llmdoc 优先**: 任何代码操作前必须先读取文档
3. **auggie 集成**: 智能代码搜索，自动上下文理解
4. **智能分流**: Phase 1 后自动判断任务复杂度
5. **Security-First**: 安全永远第一

---

## 智能分流工作流

### 自动分类逻辑

```
                    Phase 1: Discovery
                            ↓
                   智能分类 (自动判断)
                            ↓
              ┌──────────┼──────────┐
              │          │          │
              ↓          ↓          ↓
        简单任务    中等任务    复杂任务
        (Bug修复)   (用户确认)   (新功能)
              │          │          │
              │          ↓          │
              │    ┌─────────────┐   │
              │    │  询问用户    │   │
              │    │  选择流程    │   │
              │    └──────┬───────┘   │
              │           │          │
              ↓           ↓          ↓
        简单流程    简单流程    完整流程
```

### 分类信号

| 信号 | 简单任务 | 复杂任务 |
|:---|:---|:---|
| **关键词** | bug, 修复, error, fix, 问题, 调试 | 功能, 模块, 新, 开发, 重构, 系统, 设计 |
| **文件数** | <= 2 | >= 5 |
| **风险** | low | any |

---

## 8.5 阶段融合工作流

### 简单流程 (SIMPLE) - Bug 修复、小修改

| 阶段 | 名称 | Agent | 跳过 |
|:---|:---|:---|:---:|
| 0 | llmdoc + ace-tool | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 3 | Clarification | 主 Agent + AskUser | 2, 4, 4.5 |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Security + Review | `code-simplifier` + `security-guidance` + `code-reviewer` | - |
| 8 | Summary | `investigator` | - |

### 完整流程 (COMPLEX) - 新功能、模块开发

| 阶段 | 名称 | Agent | 并行 |
|:---|:---|:---|:---:|
| 0 | llmdoc + ace-tool | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 2 | Exploration | `investigator` x3 | 是 |
| 3 | Clarification | 主 Agent + AskUser | - |
| 4 | Architecture | `architect` x2 + `investigator` | 是 |
| **4.5** | **Spec** | `architect` + `ace-tool` | **-** |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Security + Review | `code-simplifier` + `security-guidance` + `code-reviewer` | 是 |
| 8 | Summary | `investigator` | - |

> 📖 **详细流程说明**: 每个阶段的具体步骤、智能分类逻辑、并行执行模式等，请参阅 [WORKFLOW.md](docs/WORKFLOW.md)

---

## 快速开始

### 方式 1: 使用插件市场安装 (推荐)

```bash
# 添加 gclm-flow 插件市场
/plugin marketplace add https://github.com/gclm/gclm-flow

# 安装 gclm 插件
/plugin install gclm@gclm-flow
```

### 方式 2: 手动安装

```bash
# 克隆插件
cd /tmp
git clone https://github.com/gclm/gclm-flow.git
cd gclm-flow

# 运行安装脚本
bash scripts/install.sh
```

### 方式 3: 直接配置

将 [`CLAUDE.example.md`](CLAUDE.example.md) 的内容复制到 `~/.claude/CLAUDE.md` 文件中。

**完成！** 配置完成后，所有行为自动激活。

### 通知功能配置（macOS 可选）

> **注意**: 通知功能仅支持 macOS

安装后，可以配置 macOS 系统通知：

```bash
# 1. 安装 terminal-notifier
brew install terminal-notifier

# 2. 安装 ClaudeNotifier.app
# 将 ClaudeNotifier.app 放到 /Applications/ 目录

# 3. 重新运行安装脚本（会自动配置 hooks）
bash scripts/install.sh
```

安装后，当 Claude Code 需要你的操作时，会收到系统通知。

---

## 工作原理

### 自动行为（无需命令）

配置 `CLAUDE.example.md` 后，这些行为**始终激活**：

| 行为 | 效果 |
|:---|:---|
| **文档优先** | Agent 在任何操作前先阅读 `llmdoc/` |
| **智能分流** | Phase 1 后自动判断任务复杂度 |
| **代码搜索** | auggie 自动提供代码上下文 |
| **选项式编程** | 不会直接下结论；通过问题呈现选择 |
| **文档维护提示** | 编码后询问是否更新文档 |
| **系统通知** | 需要操作时发送系统通知（仅 macOS） |

### 可用 Skills（自动触发）

| Skill | 触发词 | 描述 |
|:---|:---|:---|
| `/gclm` | "实现功能"、"开发新功能" | 启动智能分流工作流 |
| `/investigate` | "什么是"、"X怎么工作"、"分析" | 快速代码库调查 |
| `/tdd` | "写测试"、"TDD" | 测试驱动开发 |

---

## 融合设计

```
SpecDD (规范驱动) + TDD (测试驱动) + ace-tool (代码搜索) + llmdoc (文档优先)
```

### 关键改进

1. **SpecDD 集成**: 复杂模块先写 Spec，再写测试
2. **ace-tool 集成**: 替代 qmd，专为代码设计
3. **智能分流**: 自动判断任务类型，选择最优流程

---

## 关键约束

1. **llmdoc 优先**: Phase 0 强制执行
2. **智能分流**: Phase 1 后自动判断任务类型
3. **Phase 3 不可跳过**: 必须澄清所有疑问
4. **Phase 5 TDD 强制**: 必须先写测试
5. **并行优先**: 能并行的任务必须并行执行
6. **状态持久化**: 中途退出可恢复
7. **选项式编程**: 使用 AskUserQuestion 展示选项
8. **文档更新询问**: Phase 7 必须询问

---

## 目录结构

```
gclm-flow/
├── README.md                      # 本文件
├── CLAUDE.example.md              # 示例配置
├── .claude-plugin/                 # 插件市场配置
│   └── marketplace.json
├── agents/                        # Agent 定义
│   ├── investigator.md
│   ├── architect.md
│   ├── worker.md
│   ├── tdd-guide.md
│   └── code-reviewer.md
├── commands/                      # 命令定义
│   ├── gclm.md
│   ├── investigate.md
│   └── tdd.md
├── skills/                        # Skill 定义
│   ├── gclm/
│   │   ├── SKILL.md
│   │   └── .mcp.json
│   ├── investigate/
│   │   └── SKILL.md
│   ├── tdd-workflow/
│   │   └── SKILL.md
│   └── file-naming-helper/
│       └── SKILL.md
├── rules/                         # 规则文件
│   ├── phases.md
│   ├── llmdoc.md
│   ├── tdd.md
│   └── agents.md
├── scripts/                       # 脚本
│   └── install.sh                 # 安装脚本
├── hooks/                         # Hooks
│   ├── notify.sh                  # 通知脚本
│   └── stop-gclm-loop.sh          # 停止循环 hook
├── settings-hooks.example.json    # Hooks 配置示例
└── docs/                          # 文档
    └── WORKFLOW.md
```

---

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

---

## 外部工具集成

| 工具 | 用途 | 安装 | 必需 |
|:---|:---|:---|:---:|
| **auggie** | 代码搜索、上下文增强 | `npm install -g @augmentcode/auggie@prerelease` | ✅ |

---

## 成功指标

1. 测试通过，覆盖率 > 80%
2. 无已知安全漏洞
3. 代码可读性高
4. 需求完整满足
5. 文档已更新（如选择）

---

## 成本与效果

**诚实评估**: 智能分流工作流用 **1.5 倍的价钱**完成了从 85 分到 92 分的效果提升。

- 简单项目：快速 TDD，效率提升
- 复杂项目：SpecDD 确保 quality，收益显著
- 生产级代码库（10万+ 行）：效果出色

---

## License

MIT

---

<div align="center">

由 **gclm** 精心打造

</div>
