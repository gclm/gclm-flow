#!/usr/bin/env bash
set -euo pipefail

TMPDIR=$(mktemp -d)
cleanup() {
  rm -rf "$TMPDIR"
}
trap cleanup EXIT

run_case() {
  local name="$1"
  local prompt="$2"
  local output_file="$TMPDIR/$name.txt"
  printf '=== %s ===\n' "$name"
  codex exec --skip-git-repo-check -C "$TMPDIR" "$prompt" | tee "$output_file"
  printf '\n'
}

run_case basic 'Reply exactly HOOK_SMOKE_OK and do not use any tools.'
grep -q '^HOOK_SMOKE_OK$' "$TMPDIR/basic.txt"

run_case shell 'Use the shell tool to run exactly: echo hook_smoke_ok . After the tool finishes, reply with the exact text SHELL_HOOK_SMOKE_OK.'
grep -q '^SHELL_HOOK_SMOKE_OK$' "$TMPDIR/shell.txt"

run_case risk 'Use the shell tool to run exactly: git reset --hard . After the tool attempt finishes, explain the result in one short sentence.'
grep -Eq '风险防护钩子拦截|高风险硬重置操作' "$TMPDIR/risk.txt"

HOOK_JSON=$(mktemp "$TMPDIR/commit-ready.XXXXXX.json")
HOOK_REPO=$(mktemp -d "$TMPDIR/commit-ready-repo.XXXXXX")
git -C "$HOOK_REPO" init -q
git -C "$HOOK_REPO" config user.name smoke
git -C "$HOOK_REPO" config user.email smoke@example.com
printf 'seed
' > "$HOOK_REPO/notes.txt"
git -C "$HOOK_REPO" add notes.txt
git -C "$HOOK_REPO" commit -qm 'seed'
printf 'change
' >> "$HOOK_REPO/notes.txt"
cat > "$HOOK_JSON" <<JSON
{"cwd":"$HOOK_REPO","tool_input":{"cmd":"python3 codex/bin/lint-skills.py"},"tool_response":{"exit_code":0,"stdout":"","stderr":""}}
JSON
python3 "$HOME/.codex/hooks/post_tool_commit_ready_hint.py" < "$HOOK_JSON" > "$TMPDIR/commit-ready.txt"
grep -q 'Commit readiness' "$TMPDIR/commit-ready.txt"

echo 'HOOK_SMOKE_TEST_OK'
