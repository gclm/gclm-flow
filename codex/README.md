# Codex Configuration Source of Truth

This directory is the Git-managed source of truth for personal Codex configuration.

## Scope

Managed here:
- `config.toml`
- `AGENTS.md`
- `agents/*.toml`
- `hooks/*.py`
- `skills/`: Git-managed custom skills
- `bin/sync-to-home.sh`
- `bin/diff-home.sh`
- `bin/serve-local.sh`
- `bin/github-webhook-local.sh`
- `bin/smoke-test-hooks.sh`

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
- `skills/`: layered skill system for workflow, orchestration, and domains
- `bin/`: sync and diff scripts for publishing to `~/.codex`

## Skill Layers

- `workflow`: `brainstorming`, `writing-plans`, `systematic-debugging`, `test-driven-development`, `verification-before-completion`, `writing-skills`, `updating-domain-skills`
- `orchestration`: `using-git-worktrees`, `executing-plans`, `dispatching-parallel-agents`, `finishing-a-development-branch`
- `quality`: `code-review`, `testing`
- `domain`: `documentation`, `devops`, `database`, and language stack skills

## Workflow

1. Edit files in this directory.
2. Review with `bin/diff-home.sh`.
3. Publish with `bin/sync-to-home.sh`.
4. Keep runtime state in `~/.codex` out of Git.

## Common launchers

- `bin/serve-local.sh`: start `codex serve` with stable local defaults
- `bin/github-webhook-local.sh`: start `codex github` with env-driven webhook settings
- `bin/smoke-test-hooks.sh`: run a real hook smoke test through `codex exec`
