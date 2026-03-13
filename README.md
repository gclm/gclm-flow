# gclm-flow codex branch

This branch is a standalone source of truth for personal AI coding assistant configuration.

## Contents

- `codex/`: Codex config, hooks, agent roles, and sync scripts
- `claude/`: Claude Code config, hooks, subagents, rules, and sync scripts
- `skills/`: shared skills used by both Codex and Claude Code
- `vendor/everything-claude-code/`: ECC submodule — provides rules and skills

## Claude Code config overview

### Rules (`~/.claude/rules/`)

Loaded automatically by Claude Code — no `@` references needed in CLAUDE.md.

| Stack | Source |
|-------|--------|
| common | ECC submodule |
| golang | ECC submodule |
| python | ECC submodule |
| typescript | ECC submodule |
| java | `claude/rules/java/` (local) |
| rust | `claude/rules/rust/` (local) |

Path-scoped rules (with `paths:` frontmatter) activate only when matching files are in context.

### Skills (`~/.claude/skills/`)

Skills from `skills/` (repo root) are deployed globally. Selected ECC skills are also synced:

- `eval-harness` — eval-driven development framework
- `verification-loop` — pre-PR verification (build/types/lint/test/security/diff)
- `skill-stocktake` — periodic skill quality audit

### Commands (`~/.claude/commands/`)

Slash commands from ECC submodule:

- `/learn` — extract reusable patterns from current session
- `/skill-create` — generate skills from git history
- `/evolve` — promote instincts to skills/commands/agents
- `/instinct-status` — view learned instincts

### Continuous learning

Stop hook automatically reviews each session for reusable patterns and asks whether to record them to:
1. Global `MEMORY.md` (cross-project patterns)
2. Project memory (project-specific conventions)
3. `~/.claude/skills/learned/` (learned skills)

### Hooks

| Hook | Trigger | Purpose |
|------|---------|--------|
| `session_start_context.py` | SessionStart | Check project doc health |
| `pre_tool_risk_guard.py` | PreToolUse (Bash/Write/Edit) | Block dangerous commands, warn on sensitive files |
| `post_tool_commit_ready_hint.py` | PostToolUse (Bash) | Suggest commit when changes are ready |
| `post_tool_git_push_hint.py` | PostToolUse (Bash) | Remind to push after commit |
| `stop_self_check.py` | Stop | Verify outcome, check docs drift, continuous learning prompt |

## Publish workflow

### Codex

```bash
bash codex/bin/sync-to-home.sh
```

Publishes to `~/.codex/`. Skills are synced to `~/.agents/skills/`.

### Claude Code

```bash
# Deploy CLAUDE.md, hooks, agents, skills, rules, commands
bash claude/bin/sync-to-home.sh

# Register MCP servers (auggie, yunxiao, exa)
bash claude/bin/setup-mcp.sh
```

Publishes to `~/.claude/`. Hooks are injected into `~/.claude/settings.json`. MCP servers are registered via `claude mcp add` to user scope.

### Update ECC submodule

```bash
git submodule update --remote vendor/everything-claude-code
bash claude/bin/sync-to-home.sh
```
