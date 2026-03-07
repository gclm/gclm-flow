#!/usr/bin/env bash
set -euo pipefail

SRC_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
DST_DIR="$HOME/.codex"
BACKUP_DIR="$DST_DIR/backups/codex-config-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$DST_DIR/agents" "$DST_DIR/hooks" "$BACKUP_DIR"

backup_if_exists() {
  local rel="$1"
  local dst="$DST_DIR/$rel"
  if [[ -e "$dst" ]]; then
    mkdir -p "$BACKUP_DIR/$(dirname "$rel")"
    cp "$dst" "$BACKUP_DIR/$rel"
  fi
}

copy_managed_file() {
  local rel="$1"
  backup_if_exists "$rel"
  mkdir -p "$DST_DIR/$(dirname "$rel")"
  cp "$SRC_DIR/$rel" "$DST_DIR/$rel"
}

copy_managed_file "config.toml"
copy_managed_file "AGENTS.md"

for file in "$SRC_DIR"/agents/*.toml; do
  rel="agents/$(basename "$file")"
  copy_managed_file "$rel"
done

for file in "$SRC_DIR"/hooks/*.py; do
  rel="hooks/$(basename "$file")"
  copy_managed_file "$rel"
  chmod +x "$DST_DIR/$rel"
done

echo "Published managed Codex config to $DST_DIR"
echo "Backups stored at $BACKUP_DIR"
