# Fork-First Codex Configuration Design

## Goal

Create a Git-managed source of truth for personal Codex configuration under `gclm-flow/codex/`, with fork-first defaults, reusable agent roles, mixed-mode hooks, and a safe publish path into `~/.codex`.

## Context

The runtime directory `~/.codex` already contains authentication material, history, sessions, SQLite state, and other ephemeral data. Those files must not become the Git source of truth.

The desired model is:
- Git manages declarative configuration only.
- `~/.codex` remains the runtime directory.
- Publishing is one-way from repo to runtime.

## Chosen Structure

```text
gclm-flow/
  codex/
    config.toml
    AGENTS.md
    agents/*.toml
    hooks/*.py
    bin/sync-to-home.sh
    bin/diff-home.sh
```

## Key Decisions

### 1. Source of truth lives in `codex/`

This keeps Codex-specific configuration isolated from other gclm-flow concerns and avoids polluting the runtime directory with Git metadata.

### 2. Runtime state is excluded from Git

Files such as `auth.json`, `history.jsonl`, `state_*.sqlite*`, `sessions/`, and `shell_snapshots/` remain unmanaged.

### 3. Mixed-mode hook policy

Hooks are split into:
- blocking for clearly dangerous actions
- advisory for stability and workflow hygiene
- review-oriented reminders at session end

### 4. Agent roles are converted, not copied raw

The existing markdown role descriptions under `gclm-flow/agents/` are preserved conceptually, but converted into Codex-native role config files.

### 5. Global guidance stays principle-focused

`AGENTS.md` keeps stable operating rules, while executable enforcement moves into hooks.

## Configuration Model

Main config responsibilities:
- local provider registration
- default safety boundaries
- three profiles: `fast`, `deep`, `review`
- agent registry
- mixed-mode hook wiring
- trusted project paths
- disable persisted history writes

## Agent Mapping

High-frequency roles:
- planner
- investigator
- builder
- reviewer

Low-frequency support roles:
- recorder
- remember

## Hook Set

Initial hook set:
- `session_start_context.py`
- `pre_tool_risk_guard.py`
- `pre_tool_tmux_advice.py`
- `post_tool_git_push_hint.py`
- `stop_self_check.py`

## Publish Workflow

1. Edit files in `gclm-flow/codex/`.
2. Review differences with `codex/bin/diff-home.sh`.
3. Publish to `~/.codex` with `codex/bin/sync-to-home.sh`.
4. Keep runtime state out of Git.

## Validation Plan

- Parse TOML files successfully.
- Python hooks compile cleanly.
- Sync script is syntactically valid.
- Managed tree diff is understandable.

## Follow-up

After the design is written, implementation proceeds from the same worktree and lands on the `codex` branch.
