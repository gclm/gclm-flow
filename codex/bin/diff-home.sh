#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

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

compare_dir() {
  local rel="$1"
  local src="$SRC_DIR/$rel"
  local dst="$DST_DIR/$rel"
  echo "=== $rel ==="
  if [[ ! -e "$dst" ]]; then
    echo "missing in ~/.codex"
    return 0
  fi
  diff -ru "$dst" "$src" || true
}

report_extra_children() {
  local rel="$1"
  shift
  local dst_root="$DST_DIR/$rel"
  local src_root="$SRC_DIR/$rel"
  local entry base keep

  [[ -d "$dst_root" ]] || return 0

  for entry in "$dst_root"/*; do
    [[ -e "$entry" ]] || continue
    base="$(basename "$entry")"
    keep=0

    if [[ -e "$src_root/$base" ]]; then
      keep=1
    else
      for preserved in "$@"; do
        if [[ "$base" == "$preserved" ]]; then
          keep=1
          break
        fi
      done
    fi

    if [[ "$keep" -eq 0 ]]; then
      echo "extra in ~/.codex: $rel/$base"
    fi
  done
}

compare_file "config.toml"
compare_file "AGENTS.md"

for file in "$SRC_DIR"/agents/*.toml; do
  compare_file "agents/$(basename "$file")"
done
report_extra_children "agents"

for file in "$SRC_DIR"/hooks/*.py; do
  compare_file "hooks/$(basename "$file")"
done
report_extra_children "hooks"

for dir in "$SRC_DIR"/skills/*; do
  [[ -d "$dir" ]] || continue
  compare_dir "skills/$(basename "$dir")"
done
report_extra_children "skills" ".system"
