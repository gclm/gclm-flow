#!/usr/bin/env python3
"""Claude Code sessionStart hook.
Outputs additional context about the current working directory.
"""
import json
import os
import subprocess
import sys
from pathlib import Path


def run(cmd: list[str], cwd: str) -> str:
    result = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)
    if result.returncode != 0:
        return ""
    return result.stdout.strip()


def missing_doc_markers(path: Path) -> list[str]:
    missing = []
    if not (path / "README.md").exists():
        missing.append("README.md")
    if not (path / "CLAUDE.md").exists():
        missing.append("CLAUDE.md")
    if not (path / "docs").exists() and not (path / "llmdoc").exists():
        missing.append("docs-or-llmdoc")
    return missing


def main() -> None:
    payload = json.load(sys.stdin)
    cwd = payload.get("cwd") or os.getcwd()
    path = Path(cwd)
    facts = [f"cwd={cwd}"]
    hints = []
    if (path / "CLAUDE.md").exists():
        facts.append("project_has_claude_md=yes")
    git_root = run(["git", "rev-parse", "--show-toplevel"], cwd)
    if git_root:
        facts.append(f"git_root={git_root}")
        branch = run(["git", "branch", "--show-current"], cwd)
        if branch:
            facts.append(f"git_branch={branch}")
    for marker, label in [
        ("Cargo.toml", "rust"),
        ("go.mod", "go"),
        ("package.json", "node"),
        ("pyproject.toml", "python"),
    ]:
        if (path / marker).exists():
            facts.append(f"stack={label}")
            break
    if git_root:
        missing = missing_doc_markers(Path(git_root))
        if missing:
            hints.append(
                "Documentation bootstrap suggested: missing " + ", ".join(missing) + ". Use the documentation skill to initialize the minimal project docs."
            )
    parts = ["Session context: " + ", ".join(facts)]
    if hints:
        parts.extend(hints)
    # Claude Code hook output format
    print(json.dumps({"additionalContext": " ".join(parts)}))


if __name__ == "__main__":
    main()
