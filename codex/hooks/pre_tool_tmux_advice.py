#!/usr/bin/env python3
import json
import os
import re
import sys

LONG_RUNNING = re.compile(
    r"(npm\s+(run\s+)?dev\b|pnpm\s+(run\s+)?dev\b|yarn\s+dev\b|bun\s+run\s+dev\b|"
    r"cargo\s+test\b|go\s+test\b|pytest\b|playwright\b|docker\s+compose\s+up\b|make\b)"
)


def main() -> None:
    payload = json.load(sys.stdin)
    tool_input = payload.get("tool_input") or {}
    command = tool_input.get("command") or tool_input.get("cmd") or ""
    if os.environ.get("TMUX"):
        print("{}")
        return
    if LONG_RUNNING.search(command):
        print(json.dumps({
            "additionalContext": (
                "Long-running command detected. Consider running it inside tmux for better stability and easier recovery."
            )
        }))
        return
    print("{}")


if __name__ == "__main__":
    main()
