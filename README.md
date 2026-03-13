# gclm-flow codex branch

This branch is a standalone source of truth for personal AI coding assistant configuration.

## Contents

- `codex/`: Codex config, hooks, agent roles, and sync scripts
- `claude/`: Claude Code config, hooks, subagents, and sync scripts
- `skills/`: shared skills used by both Codex and Claude Code

## Publish workflow

### Codex

```bash
bash codex/bin/sync-to-home.sh
```

Publishes to `~/.codex/`. Skills are synced to `~/.agents/skills/`.

### Claude Code

```bash
# Deploy CLAUDE.md, hooks, agents, skills
bash claude/bin/sync-to-home.sh

# Register MCP servers (auggie, yunxiao, exa)
bash claude/bin/setup-mcp.sh
```

Publishes to `~/.claude/`. Hooks are injected into `~/.claude/settings.json`. MCP servers are registered via `claude mcp add` to user scope.
