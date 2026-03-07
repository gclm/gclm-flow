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
- `AGENTS.md`：全局运行约束
- `agents/`：可复用的 agent 角色配置
- `hooks/`：混合式门禁、提醒与审查钩子
- `skills/`：按工作流、编排层、质量层、领域层划分的 skills 体系
- `bin/`：将源码配置发布到 `~/.codex` 的同步与校验脚本
- `bin/lint-skills.py`：skills 结构校验脚本，可单独运行，也会被 hook 调用

## Skill 分层

| 分层 | 用途 | 代表 skills |
| --- | --- | --- |
| `workflow` | 定义通用工作方式，覆盖方案设计、计划编写、调试、TDD、完成前验证与经验沉淀。 | `brainstorming`、`writing-plans`、`systematic-debugging`、`test-driven-development`、`verification-before-completion`、`writing-skills`、`updating-domain-skills` |
| `orchestration` | 处理执行编排与任务推进方式，覆盖 worktree、按计划执行、并行 agent 调度与分支收尾。 | `using-git-worktrees`、`executing-plans`、`dispatching-parallel-agents`、`finishing-a-development-branch` |
| `quality` | 负责质量把关，统一代码审查和测试策略，避免各领域重复维护同类规则。 | `code-review`、`testing` |
| `domain` | 承载领域特有知识与栈差异，覆盖文档、数据库、DevOps 以及各语言框架实践。 | `documentation`、`devops`、`database`、`frontend-stack`、`python-stack`、`go-stack`、`java-stack`、`rust-stack` |

## 日常流程

1. 在当前目录修改配置。
2. 用 `bin/diff-home.sh` 检查源码与运行态差异。
3. 用 `bin/lint-skills.py` 检查 skills 结构是否漂移。
4. 用 `bin/sync-to-home.sh` 发布到 `~/.codex`。
5. 将 `~/.codex` 视为运行态副本，不直接在里面做 Git 管理。

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
- `bin/lint-skills.py`：检查 skill frontmatter、references 链接和基础结构
- `bin/smoke-test-hooks.sh`：通过 `codex exec` 跑一轮真实 hook smoke test
