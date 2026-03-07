#!/usr/bin/env python3
import json
import os
import subprocess
from pathlib import Path

DOC_EXTS = {".md", ".mdx", ".rst", ".txt"}
HIGH_RISK_PATH_PARTS = (
    "openapi",
    "schema",
    "config",
    "settings",
    "hooks",
    "agents",
    "cmd",
    "cli",
)
HIGH_RISK_EXTS = {".yaml", ".yml", ".json", ".toml", ".proto"}


def run(cmd: list[str], cwd: str) -> str:
    result = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)
    if result.returncode != 0:
        return ""
    return result.stdout.strip()


def git_status_summary(cwd: str) -> str:
    return run(["git", "status", "--short"], cwd)


def changed_files(cwd: str) -> list[str]:
    files: set[str] = set()
    diff_out = run(["git", "diff", "--name-only", "HEAD"], cwd)
    if diff_out:
        files.update(line.strip() for line in diff_out.splitlines() if line.strip())
    status_out = git_status_summary(cwd)
    if status_out:
        for line in status_out.splitlines():
            if len(line) > 3:
                files.add(line[3:].strip())
    return sorted(files)


def is_doc_file(rel: str) -> bool:
    path = Path(rel)
    if path.suffix.lower() in DOC_EXTS:
        return True
    head = path.parts[0] if path.parts else ""
    return head in {"docs", "llmdoc"}


def is_high_risk_doc_sensitive(rel: str) -> bool:
    normalized = rel.lower()
    if any(part in normalized for part in HIGH_RISK_PATH_PARTS):
        return True
    return Path(rel).suffix.lower() in HIGH_RISK_EXTS


def main() -> None:
    payload = json.load(sys.stdin)
    cwd = payload.get("cwd") or os.getcwd()
    reminders = [
        "Before ending the task, verify the requested outcome with fresh evidence.",
    ]
    status = git_status_summary(cwd)
    if status:
        reminders.append("There are local git changes. Make sure they are explained in the final response.")
    files = changed_files(cwd)
    if files:
        changed_docs = any(is_doc_file(f) for f in files)
        changed_sensitive = [f for f in files if is_high_risk_doc_sensitive(f)]
        if changed_sensitive and not changed_docs:
            reminders.append(
                "Documentation drift check: code/config/runtime-facing files changed but no docs changed. Review README/docs/llmdoc if behavior, config, commands, hooks, or API changed."
            )
    print(json.dumps({"additionalContext": " ".join(reminders)}))


if __name__ == "__main__":
    import sys
    main()
