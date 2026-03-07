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

## Maintenance Conventions

### Branch and source of truth

- Ongoing Codex configuration work happens on branch `codex`.
- `main` can remain a separately published line; do not treat it as the day-to-day maintenance branch.
- This directory is the only Git-managed source of truth. `~/.codex` is a published runtime copy.

### Layering rules

- `README.md`: maintenance contract, layout, publish workflow.
- `AGENTS.md`: runtime operating guidance for Codex.
- `hooks/`: guardrails and reminders. Only block high-risk actions; low-risk cases should warn or request review.
- `agents/*.toml`: specialized agent roles. Keep them narrow and composable.
- `skills/`: reusable judgment and workflows.

### Skill maintenance rules

- Keep `SKILL.md` thin: trigger conditions, core rules, and links.
- Move long checklists, case studies, and stack-specific details into `references/`.
- Put shared process once in global skills such as `testing`, `code-review`, and `documentation`.
- Keep domain skills focused on domain-specific deltas; do not duplicate global workflow rules into every stack.
- If two skills overlap heavily, merge or delete instead of keeping near-duplicates.
- New durable lessons should go through `updating-domain-skills`; use `agents/remember.toml` when structured extraction helps.

### `references/` writing rules

- One file, one topic.
- Start with a short purpose line.
- Prefer a stable structure: `何时查看`, `重点做法` or `重点检查`, then `注意事项` or `检查清单`.
- Record reusable decisions, verification paths, and pitfalls, not one-off narration.
- Never store secrets, personal data, machine-specific absolute paths, or temporary incident noise.

### Hooks and provider rules

- Hooks should stay small, testable, and policy-focused.
- Initialization checks, documentation drift reminders, and risk gating belong in hooks when they are cross-cutting.
- Provider, model, and web settings stay centralized in config and launcher scripts; avoid scattering them across skills.
- Changes to hooks, agents, or launchers should be validated with `bin/smoke-test-hooks.sh` or an equivalent real execution path.

### Change workflow

1. Edit under `codex/`.
2. Check `git status` and review the relevant diff.
3. Run targeted verification for the changed area.
4. Review runtime drift with `bin/diff-home.sh`.
5. Publish with `bin/sync-to-home.sh`.
6. Commit and push branch `codex`.

## Common launchers

- `bin/serve-local.sh`: start `codex serve` with stable local defaults
- `bin/github-webhook-local.sh`: start `codex github` with env-driven webhook settings
- `bin/smoke-test-hooks.sh`: run a real hook smoke test through `codex exec`
