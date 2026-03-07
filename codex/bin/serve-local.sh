#!/usr/bin/env bash
set -euo pipefail

HOST="${CODEX_SERVE_HOST:-127.0.0.1}"
PORT="${CODEX_SERVE_PORT:-8788}"
TOKEN="${CODEX_SERVE_TOKEN:-}"
DEV_FLAG="${CODEX_SERVE_DEV:-0}"

cmd=(codex serve --host "$HOST" --port "$PORT" --no-open)
if [[ -n "$TOKEN" ]]; then
  cmd+=(--token "$TOKEN")
fi
if [[ "$DEV_FLAG" == "1" ]]; then
  cmd+=(--dev)
fi

printf 'Starting Codex Web UI on http://%s:%s\n' "$HOST" "$PORT"
exec "${cmd[@]}"
