#!/usr/bin/env python3
"""Claude Code stop hook.
Before the session ends, reminds to verify outcome and check for uncommitted changes.
"""
import json
import os
import subprocess
import sys
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
DEVOPS_PATH_MARKERS = (
    "dockerfile",
    ".github/workflows",
    "k8s/",
    "helm/",
    "terraform/",
    "deploy/",
)


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


def detect_domain(cwd: str, files: list[str]) -> str | None:
    root = Path(run(["git", "rev-parse", "--show-toplevel"], cwd) or cwd)
    lowered = [f.lower() for f in files]
    if any(marker in f for f in lowered for marker in DEVOPS_PATH_MARKERS):
        return "devops"
    if any(Path(f).name.lower() == "dockerfile" for f in files):
        return "devops"
    if (root / "go.mod").exists() and any(Path(f).suffix == ".go" for f in files):
        return "go-stack"
    if (root / "Cargo.toml").exists() and any(Path(f).suffix == ".rs" for f in files):
        return "rust-stack"
    if ((root / "pyproject.toml").exists() or (root / "requirements.txt").exists()) and any(Path(f).suffix == ".py" for f in files):
        return "python-stack"
    if ((root / "pom.xml").exists() or (root / "build.gradle").exists() or (root / "build.gradle.kts").exists()) and any(Path(f).suffix in {".java", ".kt"} for f in files):
        return "java-stack"
    if (root / "package.json").exists() and any(Path(f).suffix in {".ts", ".tsx", ".js", ".jsx"} for f in files):
        return "frontend-stack"
    return None


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
        domain = detect_domain(cwd, files)
        if domain:
            reminders.append(
                f"Domain knowledge check: this task touched {domain}. If you learned a reusable pattern or pitfall, consider updating-domain-skills before wrapping up."
            )
    # Continuous learning: auto-prompt at session end when files changed
    if files:
        reminders.append(
            "CONTINUOUS LEARNING CHECK: Review this session for reusable patterns before finishing. "
            "Look for: (1) non-obvious fixes or workarounds, (2) repeated workflows worth turning into a skill, "
            "(3) project conventions discovered, (4) tool combinations that worked well. "
            "If anything is genuinely reusable, present 1-3 candidates and ask the user in Chinese: "
            "'要把这些记录到 MEMORY.md 或 learned skill 吗？' "
            "Skip this entirely if the session was trivial or only contained documentation edits."
        )
    print(json.dumps({"additionalContext": " ".join(reminders)}))


if __name__ == "__main__":
    main()
