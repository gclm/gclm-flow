#!/usr/bin/env python3
import json
import os
import subprocess
import sys


def git_status_summary(cwd: str) -> str:
    result = subprocess.run(["git", "status", "--short"], cwd=cwd, capture_output=True, text=True)
    if result.returncode != 0:
        return ""
    return result.stdout.strip()


def main() -> None:
    payload = json.load(sys.stdin)
    cwd = payload.get("cwd") or os.getcwd()
    reminders = [
        "Before ending the task, verify the requested outcome with fresh evidence.",
    ]
    status = git_status_summary(cwd)
    if status:
        reminders.append("There are local git changes. Make sure they are explained in the final response.")
    print(json.dumps({"additionalContext": " ".join(reminders)}))


if __name__ == "__main__":
    main()
