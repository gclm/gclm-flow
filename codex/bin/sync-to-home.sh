#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

SRC_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
DST_DIR="$HOME/.codex"
AGENTS_DIR="$HOME/.agents"
BACKUP_DIR="$DST_DIR/backups/codex-config-$(date +%Y%m%d-%H%M%S)"

mkdir -p "$DST_DIR/agents" "$DST_DIR/hooks" "$DST_DIR/bin" "$AGENTS_DIR/skills" "$BACKUP_DIR"

backup_if_exists() {
  local root="$1"
  local rel="$2"
  local dst="$root/$rel"
  if [[ -e "$dst" ]]; then
    mkdir -p "$BACKUP_DIR/$(dirname "$rel")"
    cp -R "$dst" "$BACKUP_DIR/$rel"
  fi
}

copy_managed_file() {
  local rel="$1"
  backup_if_exists "$DST_DIR" "$rel"
  mkdir -p "$DST_DIR/$(dirname "$rel")"
  cp "$SRC_DIR/$rel" "$DST_DIR/$rel"
}

copy_managed_dir() {
  local dst_root="$1"
  local rel="$2"
  backup_if_exists "$dst_root" "$rel"
  rm -rf "$dst_root/$rel"
  mkdir -p "$dst_root/$(dirname "$rel")"
  cp -R "$SRC_DIR/$rel" "$dst_root/$rel"
}

prune_unmanaged_children() {
  local dst_root="$1"
  local rel="$2"
  shift 2
  local src_root="$SRC_DIR/$rel"
  local entry base keep

  mkdir -p "$dst_root/$rel"

  for entry in "$dst_root/$rel"/*; do
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
      backup_if_exists "$dst_root" "$rel/$base"
      rm -rf "$dst_root/$rel/$base"
    fi
  done
}

copy_managed_file "config.toml"
copy_managed_file "AGENTS.md"

prune_unmanaged_children "$DST_DIR" "agents"
for file in "$SRC_DIR"/agents/*.toml; do
  rel="agents/$(basename "$file")"
  copy_managed_file "$rel"
done

prune_unmanaged_children "$DST_DIR" "hooks"
for file in "$SRC_DIR"/hooks/*.py; do
  rel="hooks/$(basename "$file")"
  copy_managed_file "$rel"
  chmod +x "$DST_DIR/$rel"
done

prune_unmanaged_children "$DST_DIR" "bin"
for file in "$SRC_DIR"/bin/*; do
  [[ -f "$file" ]] || continue
  rel="bin/$(basename "$file")"
  copy_managed_file "$rel"
  chmod +x "$DST_DIR/$rel"
done

# Migrate: clean up old ~/.codex/skills/ with backup
if [[ -d "$DST_DIR/skills" ]]; then
  echo "Migrating ~/.codex/skills/ → backup (now managed under ~/.agents/skills/)"
  backup_if_exists "$DST_DIR" "skills"
  rm -rf "$DST_DIR/skills"
fi

# Sync skills to ~/.agents/skills/ (source: repo root skills/)
SKILLS_SRC="$(cd "$SRC_DIR/.." && pwd)/skills"
prune_unmanaged_children "$AGENTS_DIR" "skills" ".system"
for dir in "$SKILLS_SRC"/*; do
  [[ -d "$dir" ]] || continue
  rel="skills/$(basename "$dir")"
  backup_if_exists "$AGENTS_DIR" "$rel"
  rm -rf "$AGENTS_DIR/$rel"
  mkdir -p "$AGENTS_DIR/skills"
  cp -R "$dir" "$AGENTS_DIR/$rel"
done

echo "Published managed Codex config to $DST_DIR"
echo "Skills synced to $AGENTS_DIR/skills/"
echo "Backups stored at $BACKUP_DIR"
