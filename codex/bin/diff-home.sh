#!/usr/bin/env bash
set -euo pipefail

SRC_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
DST_DIR="$HOME/.codex"

compare_file() {
  local rel="$1"
  local src="$SRC_DIR/$rel"
  local dst="$DST_DIR/$rel"
  echo "=== $rel ==="
  if [[ ! -e "$dst" ]]; then
    echo "missing in ~/.codex"
    return 0
  fi
  diff -u "$dst" "$src" || true
}

compare_file "config.toml"
compare_file "AGENTS.md"

for file in "$SRC_DIR"/agents/*.toml; do
  compare_file "agents/$(basename "$file")"
done

for file in "$SRC_DIR"/hooks/*.py; do
  compare_file "hooks/$(basename "$file")"
done
