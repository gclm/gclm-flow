#!/usr/bin/env python3
import json
import re
import subprocess
import sys
from pathlib import Path

VERIFICATION_PATTERNS = [
    r"\bpytest\b",
    r"\bgo test\b",
    r"\bcargo test\b",
    r"\bbun run build\b",
    r"\bbun run test\b",
    r"\bnpm run build\b",
    r"\bnpm test\b",
    r"\bpnpm build\b",
    r"\bpnpm test\b",
    r"\bjust test\b",
    r"\bjust fmt\b",
    r"\blint-skills\.py\b",
    r"\bdiff-home\.sh\b",
]


def run(cmd: list[str], cwd: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)


def git_root(cwd: str) -> Path | None:
    result = run(["git", "rev-parse", "--show-toplevel"], cwd)
    if result.returncode != 0:
        return None
    return Path(result.stdout.strip())


def changed_files(cwd: str) -> list[str]:
    result = run(["git", "status", "--short"], cwd)
    if result.returncode != 0:
        return []
    files = []
    for line in result.stdout.splitlines():
        if len(line) > 3:
            files.append(line[3:].strip())
    return files


def looks_like_verification(command: str) -> bool:
    lowered = command.lower()
    return any(re.search(pattern, lowered, flags=re.IGNORECASE) for pattern in VERIFICATION_PATTERNS)


def summarize_scope(files: list[str]) -> str:
    top_levels: list[str] = []
    for rel in files:
        top = Path(rel).parts[0] if Path(rel).parts else rel
        if top not in top_levels:
            top_levels.append(top)
    if not top_levels:
        return "当前改动"
    if len(top_levels) == 1:
        return top_levels[0]
    return ", ".join(top_levels[:3]) + (" 等" if len(top_levels) > 3 else "")


def main() -> None:
    payload = json.load(sys.stdin)
    tool_input = payload.get("tool_input") or {}
    tool_response = payload.get("tool_response") or {}
    cwd = payload.get("cwd") or "."
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    if not looks_like_verification(command):
        print("{}")
        return
    if isinstance(tool_response, dict):
        exit_code = tool_response.get("exit_code")
        if exit_code not in (0, None):
            print("{}")
            return
    root = git_root(cwd)
    if root is None:
        print("{}")
        return
    files = changed_files(cwd)
    if not files:
        print("{}")
        return
    scope = summarize_scope(files)
    hints = [
        "Commit readiness: verification just succeeded and local git changes remain.",
        f"Touched scope: {scope}.",
        "If the task goal is complete, proactively output a `Commit Ready` block with verification evidence, split advice, and 1-2 candidate conventional commit titles instead of waiting for the user to ask for commit again.",
    ]
    print(json.dumps({"additionalContext": "\n".join(hints)}))


if __name__ == "__main__":
    main()
