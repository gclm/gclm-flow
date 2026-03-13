# Codex 配置源码仓库

这个目录是个人 Codex 配置的 Git 管理源码目录。

## 管理范围

这里纳入 Git 管理：
- `config.toml`
- `AGENTS.md`
- `agents/*.toml`
- `hooks/*.py`
- `skills/`：自定义 skills
- `bin/sync-to-home.sh`
- `bin/diff-home.sh`
- `bin/analyze-history.py`
- `bin/lint-skills.py`
- `bin/serve-local.sh`
- `bin/github-webhook-local.sh`
- `bin/smoke-test-hooks.sh`

这里不纳入 Git 管理：
- `~/.codex/auth.json`
- `~/.codex/history.jsonl`
- `~/.codex/state_*.sqlite*`
- `~/.codex/sessions/`
- `~/.codex/archived_sessions/`
- 其他运行时状态文件

## 目录说明

- `config.toml`：Codex 主配置
- 默认 MCP：`auggie`、`exa`、`yunxiao`
- 浏览器自动化：使用 `agent-browser` skill 替代 chrome-devtools MCP
- `AGENTS.md`：全局运行约束
- `agents/`：可复用的 agent 角色配置
- `hooks/`：混合式门禁、提醒与审查钩子
- `skills/`：按工作流、编排层、质量层、领域层划分的 skills 体系
- `bin/`：将源码配置发布到 `~/.codex` 的同步与校验脚本
- `bin/analyze-history.py`：读取 `history.jsonl` 并导出结构化复盘事实
- `bin/lint-skills.py`：skills 结构校验脚本，可单独运行，也会被 hook 调用

## Skill 分层

| 分层 | 用途 | 代表 skills |
| --- | --- | --- |
| `workflow` | 定义通用工作方式，覆盖方案设计、计划编写、调试、TDD、完成前验证、历史复盘与经验沉淀。 | `brainstorming`、`writing-plans`、`systematic-debugging`、`test-driven-development`、`verification-before-completion`、`writing-skills`、`reviewing-codex-history`、`updating-domain-skills` |
| `orchestration` | 处理执行编排与任务推进方式，覆盖 worktree、按计划执行、并行 agent 调度与分支收尾。 | `using-git-worktrees`、`executing-plans`、`dispatching-parallel-agents`、`finishing-a-development-branch` |
| `quality` | 负责质量把关，统一代码审查和测试策略，避免各领域重复维护同类规则。 | `code-review`、`testing` |
| `domain` | 承载领域特有知识与栈差异，覆盖文档、数据库、DevOps 以及各语言框架实践。 | `documentation`、`devops`、`database`、`frontend-stack`、`python-stack`、`go-stack`、`java-stack`、`rust-stack` |

## Skill 清单

| Skill | 分层 | 主要用途 | 典型触发场景 |
| --- | --- | --- | --- |
| `brainstorming` | `workflow` | 在动手前澄清需求、边界与方案取舍。 | 新功能、方案设计、需求模糊的任务。 |
| `writing-plans` | `workflow` | 把需求拆成可执行计划与验收步骤。 | 多步骤实现、跨文件变更、需要显式计划。 |
| `systematic-debugging` | `workflow` | 按证据定位问题，不靠猜测修 bug。 | 线上问题、测试失败、异常行为排查。 |
| `verification-before-completion` | `workflow` | 在宣称完成前强制补齐验证证据。 | 准备说“已完成 / 已修复 / 可合并”时。 |
| `writing-skills` | `workflow` | 规范化创建或重构 skills。 | 新增 skill、拆分 references、治理 skill 结构。 |
| `reviewing-codex-history` | `workflow` | 从 `history.jsonl` 提炼工作流摩擦和经验候选。 | 做历史复盘、找高频请求、识别治理机会。 |
| `updating-domain-skills` | `workflow` | 把已验证的领域经验回写到对应 skill。 | 真实任务产出稳定经验后。 |
| `using-git-worktrees` | `orchestration` | 在独立 worktree 中隔离开发或同步工作。 | 避免污染当前工作区、做并行分支工作。 |
| `executing-plans` | `orchestration` | 严格按现有实施计划逐步落地和验证。 | 用户已给 plan，要求按步骤执行。 |
| `dispatching-parallel-agents` | `orchestration` | 用 `spawn_team` 启动蜂群模式，并行调研或实现。 | 多模块任务、需要并行调研或独立实现的场景。 |
| `finishing-a-development-branch` | `orchestration` | 处理收尾、合并、PR 或清理决策。 | 功能已实现，准备进入分支收尾阶段。 |
| `code-review` | `quality` | 做 bug、风险、回归、测试缺口审查。 | 用户要求 review、审计、验证实现质量。 |
| `testing` | `quality` | 决定该测什么、测到哪层、如何证明行为。 | 设计测试、补测试、判断验证范围。 |
| `documentation` | `domain` | 维护 README、AGENTS、文档结构与约定。 | 改文档、补说明、治理文档分层。 |
| `devops` | `domain` | 处理部署、CI/CD、镜像、发布和运行环境。 | Docker、K8s、GitHub Actions、发布链路。 |
| `database` | `domain` | 处理 schema、migration、索引和查询风险。 | SQL、SQLite、迁移、数据安全与回滚。 |
| `frontend-stack` | `domain` | 提供前端栈的项目结构、状态管理和测试差异。 | React/Vue/TSX、前端交互和构建问题。 |
| `python-stack` | `domain` | 提供 Python 服务/脚本项目的栈级建议。 | FastAPI/Flask、Python 脚本、pytest。 |
| `go-stack` | `domain` | 提供 Go 项目结构、测试和工程实践差异。 | Go 服务、CLI、`go test`、模块治理。 |
| `java-stack` | `domain` | 提供 Java/Spring/Quarkus 的栈级注意事项。 | Java 后端、构建、测试和框架约束。 |
| `rust-stack` | `domain` | 提供 Rust 项目结构、错误处理和测试约束。 | Rust/Axum/Actix、Cargo、Rust 工程问题。 |
| `gclm-commit` | `domain` | 统一提交前检查、提交信息和 `Commit Ready` 输出。 | 准备提交、需要规范 commit message。 |

## 默认 MCP

| MCP | 主要用途 | 适合场景 | 使用边界 |
| --- | --- | --- | --- |
| `auggie` | 接入 Augment/Auggie 的 MCP 能力，补充外部工作区理解或相关集成。 | 需要使用 Auggie MCP 提供的额外上下文或工作区能力时。 | 依赖外部服务可用性与该 relay 配置，不应作为本地代码搜索的默认替代。 |
| `exa` | 提供 Exa 的联网搜索 / 检索能力。 | 需要外部资料、最新网页、联网研究时。 | 属于外部检索通道，不替代本地仓库阅读；高时效或高风险事实仍应显式校验来源。 |
| `yunxiao` | 接入阿里云云效 DevOps 平台，读取工作项、流水线、代码库信息。 | 需要查询云效工作项、发布流水线、代码评审等平台数据时。 | 依赖 `YUNXIAO_ACCESS_TOKEN` 环境变量，只用于云效平台上下文，不替代本地代码搜索。 |

> 浏览器自动化改用 `agent-browser` skill，不再使用 `chrome-devtools` MCP，避免后台残留 Chrome 进程。

### 默认 MCP 使用约定

- 先做本地搜索和代码阅读，再决定是否调用 `auggie` 或 `exa`。
- 浏览器端问题（DOM、控制台、网络请求）使用 `agent-browser` skill。
- `exa` 更适合外部事实补充，`auggie` 更适合补充外部工作区或集成侧能力；两者都不应替代本地代码理解。
- 通过 MCP 拿到证据后，输出里应说明用了哪个 MCP，以及它解决了什么本地上下文无法直接证明的问题。
- 增加新的默认 MCP 时，README 需要同时补“用途”和“使用边界”，避免默认能力变成隐性依赖。

## 首次初始化

### 1. 环境变量配置

在 `~/.zshrc` 里添加以下变量（替换为真实值）：

```bash
# Auggie MCP
export AUGMENT_API_TOKEN="your_augment_api_token"
export AUGMENT_API_URL="https://acemcp.your-relay.com/relay/"

# 云效 MCP
export YUNXIAO_ACCESS_TOKEN="your_yunxiao_access_token"
```

然后 `source ~/.zshrc` 使其生效。

### 2. MCP 命令软链接

MCP server 命令使用软链接指向 fnm 管理的 node 版本，避免后台残留多个 npm 进程：

```bash
# 查看当前 node 版本
fnm current

# 创建软链接（替换 v22.19.0 为实际版本）
ln -sf ~/.local/share/fnm/node-versions/v22.19.0/installation/bin/auggie /usr/local/bin/auggie
ln -sf ~/.local/share/fnm/node-versions/v22.19.0/installation/bin/alibabacloud-devops-mcp-server /usr/local/bin/yunxiao-mcp

# 验证
which auggie
which yunxiao-mcp
```

> 切换 node 版本后需要重新创建软链接。

### 3. 发布配置到运行目录

```bash
bash bin/sync-to-home.sh
```

## 日常流程

1. 在当前目录修改配置。
2. 用 `bin/diff-home.sh` 检查源码与运行态差异。
3. 用 `bin/lint-skills.py` 检查 skills 结构是否漂移。
4. 用 `bin/sync-to-home.sh` 发布到 `~/.codex`。
5. 将 `~/.codex` 视为运行态副本，不直接在里面做 Git 管理。

## 复盘入口

- `python3 ~/.codex/bin/analyze-history.py --output /tmp/codex-history.json`：从 `~/.codex/history.jsonl` 导出结构化事实
- `python3 ~/.codex/bin/analyze-history.py --topic devops --topic-kind domain --topic-samples 12 --format markdown`：对某个主题做定向抽样
- `reviewing-codex-history`：基于 JSON 结果做候选经验簇命名、分类和路由建议
- 复盘结论先分流：workflow/config/tooling 改进走全局治理，领域经验再交给 `updating-domain-skills`

## 维护约定

### 分支与事实来源

- 日常 Codex 配置维护统一在 `codex` 分支进行。
- `main` 可以保留为独立发布线，不作为日常维护主线。
- 当前目录是唯一 Git 管理事实来源，`~/.codex` 只是发布后的运行态副本。

### 分层职责

- `README.md`：维护者约定、结构说明、发布流程。
- `AGENTS.md`：给模型读取的运行规则与行为边界。
- `hooks/`：负责跨领域门禁、提醒与审查；仅阻断高风险动作。
- `agents/*.toml`：职责单一、可组合的专用代理。
- `skills/`：承载可复用判断、流程与领域经验。

### 语言约定

- `README.md` 优先使用中文，面向配置维护者。
- `AGENTS.md` 优先使用英文，面向模型运行时读取。
- 用户交互保持中文；命令、代码、配置键名保持英文。

### Skill 维护规则

- `SKILL.md` 保持薄入口，只写触发条件、核心规则和 references 链接。
- 长清单、案例、框架细节、排障经验统一拆到 `references/`。
- `testing`、`code-review`、`documentation` 这类全局流程只保留一份，不在各语言栈里重复。
- 领域 skill 只补充领域特有差异，不复制通用工作流。
- 两个 skill 如果高度重叠，应合并或删除，而不是长期并存。
- 新的稳定经验通过 `updating-domain-skills` 回写；需要结构化提炼时优先使用 `agents/remember.toml`。

### `references/` 编写规则

- 一文件一主题。
- 开头先写一句用途说明。
- 结构尽量稳定：`何时查看`、`重点做法` 或 `重点检查`、`注意事项` 或 `检查清单`。
- 记录可复用决策、验证路径、常见坑，不写一次性流水账。
- 不写 secrets、个人隐私、机器专属绝对路径、临时事故噪音。

### Hooks 与 provider 规则

- hooks 保持短小、可测试、以策略约束为主。
- 初始化检测、文档漂移提醒、风险门禁这类跨领域逻辑应落到 hooks。
- skill 结构类明显漂移由 hook 调用 `bin/lint-skills.py` 自动拦截。
- 提交已成为自然下一步时，通过轻量 reminder 主动给出 `Commit Ready`，不要等用户重复口述 commit。
- provider、model、web 相关配置集中放在主配置和启动脚本，不分散到 skills。
- 修改 hooks、agents、launchers 后，至少跑一条真实执行链路，例如 `bin/smoke-test-hooks.sh`。

### 变更流程

1. 在 `codex/` 下修改源码配置。
2. 先看 `git status`，再看相关 diff。
3. 先跑 `bin/lint-skills.py` 或对应专项验证。
4. 用 `bin/diff-home.sh` 检查运行态漂移。
5. 用 `bin/sync-to-home.sh` 发布到 `~/.codex`。
6. 提交并推送 `codex` 分支。

## 常用启动脚本

- `bin/serve-local.sh`：用本地稳定参数启动 `codex serve`
- `bin/github-webhook-local.sh`：用环境变量方式启动 `codex github`
- `bin/analyze-history.py`：基于 `history.jsonl` 导出可供 skill 和模型继续梳理的结构化事实
- `bin/lint-skills.py`：检查 skill frontmatter、references 链接和基础结构
- `bin/smoke-test-hooks.sh`：通过 `codex exec` 跑一轮真实 hook smoke test
