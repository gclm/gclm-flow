#!/usr/bin/env bash
# Diff claude/ source against ~/.claude/ deployed config
# Usage: bash claude/bin/diff-home.sh
set -euo pipefail

SRC="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$HOME/.claude"

diff_file() {
  local src="$1" dest="$2"
  if [ ! -f "$dest" ]; then
    echo "[MISSING] $dest"
    return
  fi
  if ! diff -q "$src" "$dest" > /dev/null 2>&1; then
    echo "[DIFF] $dest"
    diff "$src" "$dest" || true
  fi
}

diff_file "$SRC/CLAUDE.md" "$DEST/CLAUDE.md"
diff_file "$SRC/hooks.json" "$DEST/hooks.json"

for f in "$SRC/hooks/"*.py; do
  diff_file "$f" "$DEST/hooks/$(basename "$f")"
done

for skill_dir in "$SRC/skills/"/*/; do
  name="$(basename "$skill_dir")"
  for f in "$skill_dir"**/*; do
    [ -f "$f" ] || continue
    rel="${f#$skill_dir}"
    diff_file "$f" "$DEST/skills/$name/$rel"
  done
done

echo "Diff complete."
