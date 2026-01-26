# gclm-flow - 融合开发工作流插件

> 融合 myclaude/do、cc-plugin/Context Floor、everything-claude-code/TDD 的最佳实践

<div align="center">

**TDD-First + llmdoc 优先 + 多 Agent 并行**

</div>

---

## 核心哲学

1. **TDD-First**: 坚持测试驱动开发，实现代码前先写测试
2. **llmdoc 优先**: 任何代码操作前必须先读取文档
3. **需求先讨论**: Discovery → Exploration → Clarification
4. **Parallel Execution**: 尽可能并行执行任务
5. **Security-First**: 安全永远第一

---

## 8 阶段融合工作流

| 阶段 | 名称 | Agent | 并行 |
|:---|:---|:---|:---:|
| 0 | llmdoc 优先读取 | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 2 | Exploration | `investigator` x3 | 是 |
| 3 | Clarification | 主 Agent + AskUser | - |
| 4 | Architecture | `architect` x2 + `investigator` | 是 |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Doc | `worker` + `code-reviewer` | 是 |
| 8 | Summary | `investigator` | - |

> 📖 **详细流程说明**: 每个阶段的具体步骤、状态管理、并行执行模式等，请参阅 [WORKFLOW.md](docs/WORKFLOW.md)

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

---

## 工作原理

### 自动行为（无需命令）

配置 `CLAUDE.example.md` 后，这些行为**始终激活**：

| 行为 | 效果 |
|:---|:---|
| **文档优先** | Agent 在任何操作前先阅读 `llmdoc/` |
| **智能调研** | 使用 `investigator` agent 而非通用探索 |
| **选项式编程** | 不会直接下结论；通过问题呈现选择 |
| **文档维护提示** | 编码后询问是否更新文档 |

### 可用 Skills（自动触发）

| Skill | 触发词 | 描述 |
|:---|:---|:---|
| `/gclm` | "实现功能"、"开发新功能" | 启动完整工作流 |
| `/investigate` | "什么是"、"X怎么工作"、"分析" | 快速代码库调查 |
| `/tdd` | "写测试"、"TDD" | 测试驱动开发 |

---

## 融合设计

```
myclaude/do (7阶段框架) + cc-plugin/llmdoc (Context 解决方案) + everything/TDD (实现方法)
```

### 来源分析

| 来源 | 核心特点 | 约束 |
|:---|:---|:---|
| **myclaude/do** | 7阶段结构化流程、状态持久化、强制澄清 | Phase 3 不可跳过、Phase 5 需审批 |
| **cc-plugin/Context Floor** | llmdoc 文档系统、SubAgent RAG | 文档优先读取 |
| **everything-claude-code/TDD** | Red-Green-Refactor、80%覆盖率 | 测试必须先失败 |

---

## 关键约束

1. **llmdoc 优先**: Phase 0 强制执行
2. **TDD 强制**: Phase 5 必须先写测试
3. **Phase 3 不可跳过**: 必须澄清所有疑问
4. **并行优先**: 能并行的任务必须并行执行
5. **状态持久化**: 中途退出可恢复
6. **选项式编程**: 使用 AskUserQuestion 展示选项
7. **文档更新询问**: Phase 7 必须询问

---

## 目录结构

```
gclm-flow/
├── README.md                      # 本文件
├── CLAUDE.example.md              # 示例配置
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
│   │   └── SKILL.md
│   ├── investigate/
│   │   └── SKILL.md
│   └── tdd-workflow/
│       └── SKILL.md
├── rules/                         # 规则文件
│   ├── phases.md
│   ├── llmdoc.md
│   ├── tdd.md
│   └── agents.md
├── scripts/                       # 脚本
│   └── setup-gclm.sh
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
| `code-reviewer` | 代码审查 | Sonnet 4.5 |

---

## 成功指标

1. 测试通过，覆盖率 > 80%
2. 无已知安全漏洞
3. 代码可读性高
4. 需求完整满足
5. 文档已更新（如选择）

---

## 成本与效果

**诚实评估**: 这套方案大概用 **1.5 倍的价钱**完成了从 85 分到 90 分的效果提升。

- 简单项目：效果一般
- 复杂项目：收益显著
- 生产级代码库（10万+ 行）：效果出色

---

## License

MIT

---

<div align="center">

由 **gclm** 精心打造

</div>
