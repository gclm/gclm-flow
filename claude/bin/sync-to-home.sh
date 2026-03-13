#!/usr/bin/env bash
# Sync claude/ config to ~/.claude/
# Usage: bash claude/bin/sync-to-home.sh
set -euo pipefail

SRC="$(cd "$(dirname "$0")/.." && pwd)"
DEST="$HOME/.claude"

echo "Syncing $SRC -> $DEST"

# CLAUDE.md
cp "$SRC/CLAUDE.md" "$DEST/CLAUDE.md"
echo "  CLAUDE.md"

# hooks -> settings.json
python3 "$SRC/bin/inject-hooks.py" "$SRC/hooks.json" "$DEST/settings.json"
echo "  hooks.json -> settings.json"

# hooks/
mkdir -p "$DEST/hooks"
for f in "$SRC/hooks/"*.py; do
  cp "$f" "$DEST/hooks/"
  echo "  hooks/$(basename "$f")"
done

# skills/ -> ~/.claude/skills/ (source: repo root skills/)
SKILLS_SRC="$(cd "$SRC/.." && pwd)/skills"
if [ -d "$SKILLS_SRC" ]; then
  mkdir -p "$DEST/skills"
  for skill_dir in "$SKILLS_SRC/"/*/; do
    name="$(basename "$skill_dir")"
    rm -rf "$DEST/skills/$name"
    cp -r "$skill_dir" "$DEST/skills/$name"
    echo "  skills/$name"
  done
else
  echo "  [skip] no repo root skills/ found"
fi

# agents/ -> ~/.claude/agents/
if [ -d "$SRC/agents" ]; then
  mkdir -p "$DEST/agents"
  for f in "$SRC/agents/"*.md; do
    cp "$f" "$DEST/agents/"
    echo "  agents/$(basename "$f")"
  done
fi

echo "Done."
echo ""
echo "Note: MCP servers are managed separately. Run:"
echo "  bash $SRC/bin/setup-mcp.sh"
