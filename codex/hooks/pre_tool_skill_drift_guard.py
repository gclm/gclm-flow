#!/usr/bin/env python3
import json
import subprocess
import sys
from pathlib import Path

INTERESTING_COMMANDS = (
    "git commit",
    "git push",
    "sync-to-home.sh",
)


def run(cmd: list[str], cwd: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)


def git_root(cwd: str) -> Path | None:
    result = run(["git", "rev-parse", "--show-toplevel"], cwd)
    if result.returncode != 0:
        return None
    return Path(result.stdout.strip())


def changed_skill_files(root: Path) -> list[str]:
    result = run(["git", "status", "--short", "--untracked-files=all"], str(root))
    if result.returncode != 0:
        return []
    changed = []
    for line in result.stdout.splitlines():
        if len(line) <= 3:
            continue
        path = line[3:].strip()
        if path.startswith("codex/skills/"):
            changed.append(path)
    return changed


def main() -> None:
    payload = json.load(sys.stdin)
    tool_input = payload.get("tool_input") or {}
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    if not any(token in command for token in INTERESTING_COMMANDS):
        print("{}")
        return
    cwd = payload.get("cwd") or "."
    root = git_root(cwd)
    if root is None or not (root / "codex" / "skills").exists():
        print("{}")
        return
    if not changed_skill_files(root):
        print("{}")
        return
    lint_script = root / "codex" / "bin" / "lint-skills.py"
    if not lint_script.exists():
        print("{}")
        return
    result = run([sys.executable, str(lint_script)], str(root))
    if result.returncode == 0:
        print("{}")
        return
    details = (result.stdout + result.stderr).strip()
    print(
        "[codex-skill-drift-guard] blocked command because skill lint failed. "
        "Run `python3 codex/bin/lint-skills.py` and fix the reported drift before commit/push/sync.\n"
        + details,
        file=sys.stderr,
    )
    sys.exit(2)


if __name__ == "__main__":
    main()
