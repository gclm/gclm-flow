# Codex Configuration Source of Truth

This directory is the Git-managed source of truth for personal Codex configuration.

## Scope

Managed here:
- `config.toml`
- `AGENTS.md`
- `agents/*.toml`
- `hooks/*.py`
- `bin/sync-to-home.sh`
- `bin/diff-home.sh`

Not managed here:
- `~/.codex/auth.json`
- `~/.codex/history.jsonl`
- `~/.codex/state_*.sqlite*`
- `~/.codex/sessions/`
- `~/.codex/archived_sessions/`
- other runtime state

## Layout

- `config.toml`: main Codex config
- `AGENTS.md`: global operating guidance
- `agents/`: reusable agent role configs
- `hooks/`: mixed-mode guardrails and reminders
- `bin/`: sync and diff scripts for publishing to `~/.codex`

## Workflow

1. Edit files in this directory.
2. Review with `bin/diff-home.sh`.
3. Publish with `bin/sync-to-home.sh`.
4. Keep runtime state in `~/.codex` out of Git.
