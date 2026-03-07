#!/usr/bin/env python3
import json
import re
import sys

PR_URL = re.compile(r"https://github\.com/[^/]+/[^/]+/pull/\d+")


def main() -> None:
    payload = json.load(sys.stdin)
    tool_input = payload.get("tool_input") or {}
    tool_response = payload.get("tool_response") or {}
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    output = ""
    if isinstance(tool_response, dict):
        output = tool_response.get("output") or tool_response.get("stdout") or ""
    hints = []
    if "git push" in command:
        hints.append("Push completed. Confirm review status and verification evidence before treating the work as done.")
    if "gh pr create" in command:
        match = PR_URL.search(output)
        if match:
            hints.append(f"PR created: {match.group(0)}")
        hints.append("Run a review pass and confirm CI/watch status if this branch is meant to be merged.")
    if hints:
        print(json.dumps({"additionalContext": "\n".join(hints)}))
        return
    print("{}")


if __name__ == "__main__":
    main()
