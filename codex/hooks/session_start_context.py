#!/usr/bin/env python3
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


def main() -> None:
    payload = json.load(sys.stdin)
    cwd = payload.get("cwd") or os.getcwd()
    path = Path(cwd)
    facts = []
    facts.append(f"cwd={cwd}")
    if (path / "AGENTS.md").exists():
        facts.append("project_has_agents_md=yes")
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
    print(json.dumps({"additionalContext": "Session context: " + ", ".join(facts)}))


if __name__ == "__main__":
    main()
