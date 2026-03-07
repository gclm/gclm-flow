#!/usr/bin/env bash
set -euo pipefail

: "${GITHUB_WEBHOOK_SECRET:?GITHUB_WEBHOOK_SECRET is required}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN is required}"

LISTEN="${CODEX_GITHUB_LISTEN:-127.0.0.1:8787}"
MIN_PERMISSION="${CODEX_GITHUB_MIN_PERMISSION:-triage}"
COMMAND_PREFIX="${CODEX_GITHUB_COMMAND_PREFIX:-/codex}"
DELIVERY_TTL_DAYS="${CODEX_GITHUB_DELIVERY_TTL_DAYS:-7}"
REPO_TTL_DAYS="${CODEX_GITHUB_REPO_TTL_DAYS:-0}"
ALLOW_REPO="${CODEX_GITHUB_ALLOW_REPO:-gclm/gclm-flow}"

cmd=(
  codex github
  --listen "$LISTEN"
  --min-permission "$MIN_PERMISSION"
  --command-prefix "$COMMAND_PREFIX"
  --delivery-ttl-days "$DELIVERY_TTL_DAYS"
  --repo-ttl-days "$REPO_TTL_DAYS"
)

if [[ -n "$ALLOW_REPO" ]]; then
  cmd+=(--allow-repo "$ALLOW_REPO")
fi

printf 'Starting codex github on http://%s for repo %s\n' "$LISTEN" "$ALLOW_REPO"
exec "${cmd[@]}"
