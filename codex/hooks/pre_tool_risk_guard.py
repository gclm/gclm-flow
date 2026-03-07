#!/usr/bin/env python3
import json
import re
import sys
from pathlib import PurePath

SENSITIVE_PATH_MARKERS = [
    "/etc/",
    "/System/",
    "/usr/",
    "/bin/",
    "/sbin/",
    "~/.ssh/",
    ".ssh/",
    ".gnupg/",
    ".aws/",
    ".kube/",
]

DANGEROUS_COMMAND_PATTERNS = [
    (re.compile(r"(^|\s)rm\s+-rf\s+(/|~|\.|\*|/\*|~/\*|\.\./)"), "dangerous recursive delete"),
    (re.compile(r"git\s+reset\s+--hard"), "hard reset requires explicit alignment"),
    (re.compile(r"git\s+push(?:\s+[^\n]*)?\s+--force(?:-with-lease)?\b|git\s+push\s+--force(?:-with-lease)?\b"), "force push is blocked"),
    (re.compile(r"(^|\s)(shutdown|reboot|halt|poweroff)\b"), "system power command is blocked"),
    (re.compile(r"(^|\s)(chmod|chown)\b"), "permission-changing command is blocked"),
]

WRITE_TOOL_NAMES = {"write", "edit", "multi_edit", "apply_patch"}


def emit_block(reason: str) -> None:
    print(reason, file=sys.stderr)
    sys.exit(2)


def looks_sensitive_path(path: str) -> bool:
    normalized = path.replace("\\", "/")
    return any(marker in normalized for marker in SENSITIVE_PATH_MARKERS)


def check_shell(payload: dict) -> None:
    tool_input = payload.get("tool_input") or {}
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    for pattern, reason in DANGEROUS_COMMAND_PATTERNS:
        if pattern.search(command):
            emit_block(f"[codex-risk-guard] blocked shell command: {reason}")


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
            emit_block(f"[codex-risk-guard] blocked write to sensitive path: {candidate}")


def main() -> None:
    payload = json.load(sys.stdin)
    tool_name = payload.get("tool_name", "")
    if tool_name in {"shell", "exec_command"}:
        check_shell(payload)
    if tool_name in WRITE_TOOL_NAMES:
        check_write(payload)
    print("{}")


if __name__ == "__main__":
    main()
