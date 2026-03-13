# claude/

Claude Code 配置源码，对应 `codex/` 的 Claude Code 版本。
部署目标：`~/.claude/`

## 目录结构

```
claude/
├── CLAUDE.md          # 全局规则（对应 codex/AGENTS.md）
├── hooks.json         # Hooks 配置
├── hooks/             # Hook 脚本
│   ├── session_start_context.py
│   ├── pre_tool_risk_guard.py
│   ├── post_tool_commit_ready_hint.py
│   ├── post_tool_git_push_hint.py
│   └── stop_self_check.py
├── agents/            # Subagent 定义（对应 codex/agents/*.toml）
│   ├── planner.md
│   ├── investigator.md
│   ├── builder.md
│   ├── reviewer.md
│   ├── recorder.md
│   └── remember.md
└── bin/
    ├── sync-to-home.sh   # 发布到 ~/.claude/
    ├── setup-mcp.sh      # 注册 MCP servers（user scope → ~/.claude.json）
    └── diff-home.sh      # 对比源码与运行态差异
```

## Skills

Skills 格式与 codex 完全兼容，直接复用 `codex/skills/`。
`sync-to-home.sh` 会将 `claude/skills/` 同步到 `~/.claude/skills/`。

如需独立维护，将 `codex/skills/` 中的目录复制到 `claude/skills/` 即可。

## 部署

```bash
bash claude/bin/sync-to-home.sh
```

## 与 codex/ 的主要差异

| 项目 | codex/ | claude/ |
|---|---|---|
| 规则文件 | `AGENTS.md` | `CLAUDE.md` |
| Hook 配置 | `config.toml [[hooks]]` | `hooks.json` |
| Tool 名称 | `shell`, `exec_command`, `write`, `edit` | `Bash`, `Shell`, `Write`, `Edit` |
| Multi-Agent | `spawn_team` | `Agent tool` + `subagent_type` |
| 模型配置 | `config.toml` | Cursor Settings UI |
| Agent profiles | `agents/*.toml` | `CLAUDE.md` 约定 |
