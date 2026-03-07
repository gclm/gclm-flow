# gclm-flow codex branch

This branch is a standalone source of truth for personal Codex configuration.

## Contents

- `codex/`: managed Codex config, hooks, agent roles, and sync scripts
- `docs/plans/`: design and implementation notes for the Codex setup

## Publish workflow

From this branch, use `codex/bin/sync-to-home.sh` to publish managed files into `~/.codex`.
