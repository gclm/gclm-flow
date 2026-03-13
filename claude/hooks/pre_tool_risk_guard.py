#!/usr/bin/env python3
"""Claude Code preToolUse hook.
Blocks dangerous shell commands and writes to sensitive paths.
Warns on sensitive file reads and dev server without tmux.
Exit 2 = block the tool call.
"""
import json
import os
import re
import sys
from pathlib import PurePath

# Only match absolute system-level paths to avoid false positives on project dirs like claude/bin/
SENSITIVE_PATH_MARKERS = [
    "/etc/",
    "/System/",
    "/usr/",
    "/sbin/",
    "~/.ssh/",
    ".ssh/",
    ".gnupg/",
    ".aws/",
    ".kube/",
]

# Paths that must start with these prefixes (absolute system bin dirs)
SENSITIVE_PATH_PREFIXES = [
    "/bin/",
    "/sbin/",
]

# File patterns that warrant a warning when read
SENSITIVE_FILE_PATTERN = re.compile(
    r"\.(env|key|pem)$|\.env\.|credentials|secret", re.IGNORECASE
)

DANGEROUS_COMMAND_PATTERNS = [
    (re.compile(r"(^|\s)rm\s+-rf\s+(/|~|\.|\*|/\*|~/\*|\.\./)"), "dangerous recursive delete"),
    (re.compile(r"git\s+reset\s+--hard"), "hard reset requires explicit alignment"),
    (re.compile(r"git\s+push(?:\s+[^\n]*)?\s+--force(?:-with-lease)?\b|git\s+push\s+--force(?:-with-lease)?\b"), "force push is blocked"),
    (re.compile(r"(^|\s)(shutdown|reboot|halt|poweroff)\b"), "system power command is blocked"),
    (re.compile(r"(^|\s)(chmod|chown)\b"), "permission-changing command is blocked"),
]

# Dev server commands that should run inside tmux
DEV_SERVER_PATTERN = re.compile(
    r"\b(npm\s+run\s+dev|pnpm(\s+run)?\s+dev|yarn\s+dev|bun\s+run\s+dev)\b"
)
TMUX_LAUNCHER_PATTERN = re.compile(
    r"^\s*tmux\s+(new|new-session|new-window|split-window)\b"
)

WRITE_TOOL_NAMES = {"Write", "Edit", "MultiEdit"}
READ_TOOL_NAMES = {"Read", "View"}


def emit_block(reason: str) -> None:
    print(reason, file=sys.stderr)
    sys.exit(2)


def emit_warn(reason: str) -> None:
    print(reason, file=sys.stderr)


def looks_sensitive_path(path: str) -> bool:
    normalized = path.replace("\\", "/")
    if any(marker in normalized for marker in SENSITIVE_PATH_MARKERS):
        return True
    return any(normalized.startswith(prefix) for prefix in SENSITIVE_PATH_PREFIXES)


def check_shell(payload: dict) -> None:
    tool_input = payload.get("tool_input") or {}
    command = tool_input.get("command") or tool_input.get("cmd") or ""

    # Block dangerous patterns
    for pattern, reason in DANGEROUS_COMMAND_PATTERNS:
        if pattern.search(command):
            emit_block(f"[risk-guard] blocked shell command: {reason}")

    # Warn: dev server should run in tmux
    if DEV_SERVER_PATTERN.search(command) and not TMUX_LAUNCHER_PATTERN.search(command):
        if not os.environ.get("TMUX"):
            emit_warn(
                "[risk-guard] WARNING: dev server should run inside tmux for log access.\n"
                "  Suggested: tmux new-session -d -s dev '<your dev command>'"
            )


def extract_candidate_paths(tool_input: dict) -> list[str]:
    candidates = []
    for key in ("file_path", "path", "target_file", "destination"):
        value = tool_input.get(key)
        if isinstance(value, str):
            candidates.append(value)
    return candidates


def check_write(payload: dict) -> None:
    tool_input = payload.get("tool_input") or {}
    for candidate in extract_candidate_paths(tool_input):
        try:
            path = str(PurePath(candidate))
        except Exception:
            path = candidate
        if looks_sensitive_path(path):
            emit_block(f"[risk-guard] blocked write to sensitive path: {candidate}")


def check_read(payload: dict) -> None:
    tool_input = payload.get("tool_input") or {}
    for candidate in extract_candidate_paths(tool_input):
        if SENSITIVE_FILE_PATTERN.search(candidate):
            emit_warn(
                f"[risk-guard] WARNING: reading sensitive file: {candidate}\n"
                "  Ensure this data is not exposed in outputs or logs."
            )


def main() -> None:
    payload = json.load(sys.stdin)
    tool_name = payload.get("tool_name", "")
    if tool_name in {"Bash", "Shell"}:
        check_shell(payload)
    if tool_name in WRITE_TOOL_NAMES:
        check_write(payload)
    if tool_name in READ_TOOL_NAMES:
        check_read(payload)
    print("{}")


if __name__ == "__main__":
    main()
