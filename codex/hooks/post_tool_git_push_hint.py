#!/usr/bin/env python3
import json
import re
import subprocess
import sys
from pathlib import Path

PR_URL = re.compile(r"https://github\.com/[^/]+/[^/]+/pull/\d+")
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


def changed_files(cwd: str) -> list[str]:
    out = run(["git", "diff", "--name-only", "HEAD~1", "HEAD"], cwd)
    if not out:
        out = run(["git", "status", "--short"], cwd)
        if out:
            return [line[3:].strip() for line in out.splitlines() if len(line) > 3]
        return []
    return [line.strip() for line in out.splitlines() if line.strip()]


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
    tool_input = payload.get("tool_input") or {}
    tool_response = payload.get("tool_response") or {}
    cwd = payload.get("cwd") or "."
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    output = ""
    if isinstance(tool_response, dict):
        output = tool_response.get("output") or tool_response.get("stdout") or ""
    hints = []
    domain = detect_domain(cwd, changed_files(cwd))
    if "git push" in command:
        hints.append("Push completed. Confirm review status and verification evidence before treating the work as done.")
        if domain:
            hints.append(f"This push included {domain} work. If it produced reusable lessons, consider running updating-domain-skills.")
    if "gh pr create" in command:
        match = PR_URL.search(output)
        if match:
            hints.append(f"PR created: {match.group(0)}")
        hints.append("Run a review pass and confirm CI/watch status if this branch is meant to be merged.")
        if domain:
            hints.append(f"If this PR captured reusable {domain} lessons, consider updating-domain-skills after the mergeable version is settled.")
    if hints:
        print(json.dumps({"additionalContext": "\n".join(hints)}))
        return
    print("{}")


if __name__ == "__main__":
    main()
