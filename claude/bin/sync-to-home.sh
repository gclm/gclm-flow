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

# ECC commands/ -> ~/.claude/commands/ (selected from submodule)
ECC_COMMANDS="$(cd "$SRC/.." && pwd)/vendor/everything-claude-code/commands"
ECC_COMMAND_LIST="learn.md skill-create.md evolve.md instinct-status.md"
if [ -d "$ECC_COMMANDS" ]; then
  mkdir -p "$DEST/commands"
  for name in $ECC_COMMAND_LIST; do
    if [ -f "$ECC_COMMANDS/$name" ]; then
      cp "$ECC_COMMANDS/$name" "$DEST/commands/$name"
      echo "  commands/$name (from ECC)"
    fi
  done
else
  echo "  [skip] ECC submodule not found for commands"
fi

# ECC skills/ -> ~/.claude/skills/ (selected from submodule)
ECC_SKILLS="$(cd "$SRC/.." && pwd)/vendor/everything-claude-code/skills"
ECC_SKILL_LIST="eval-harness verification-loop skill-stocktake"
if [ -d "$ECC_SKILLS" ]; then
  mkdir -p "$DEST/skills"
  for name in $ECC_SKILL_LIST; do
    if [ -d "$ECC_SKILLS/$name" ]; then
      rm -rf "$DEST/skills/$name"
      cp -r "$ECC_SKILLS/$name" "$DEST/skills/$name"
      echo "  skills/$name (from ECC)"
    fi
  done
else
  echo "  [skip] ECC submodule not found for skills"
fi

# rules/ -> ~/.claude/rules/
# common + selected stacks from ECC submodule, java + rust from local rules/
ECC_RULES="$(cd "$SRC/.." && pwd)/vendor/everything-claude-code/rules"
LOCAL_RULES="$SRC/rules"
if [ -d "$ECC_RULES" ]; then
  mkdir -p "$DEST/rules"
  for stack in common golang python typescript; do
    if [ -d "$ECC_RULES/$stack" ]; then
      rm -rf "$DEST/rules/$stack"
      cp -r "$ECC_RULES/$stack" "$DEST/rules/$stack"
      echo "  rules/$stack (from ECC)"
    fi
  done
else
  echo "  [skip] vendor/everything-claude-code not found, run: git submodule update --init"
fi
for stack in java rust; do
  if [ -d "$LOCAL_RULES/$stack" ]; then
    rm -rf "$DEST/rules/$stack"
    cp -r "$LOCAL_RULES/$stack" "$DEST/rules/$stack"
    echo "  rules/$stack (local)"
  fi
done

echo "Done."
echo ""
echo "Note: MCP servers are managed separately. Run:"
echo "  bash $SRC/bin/setup-mcp.sh"
