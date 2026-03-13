#!/usr/bin/env bash
# Register MCP servers for Claude Code (user scope → ~/.claude.json)
# Usage: bash claude/bin/setup-mcp.sh
# Re-running is safe: existing servers are removed and re-added.
# Note: uses add-json to avoid issues with = in env variable values.
set -euo pipefail

remove_if_exists() {
  local name="$1"
  if claude mcp get "$name" &>/dev/null; then
    claude mcp remove --scope user "$name"
    echo "  removed existing: $name"
  fi
}

echo "Setting up MCP servers (user scope)..."

# auggie
remove_if_exists auggie
claude mcp add-json --scope user auggie \
  '{"type":"stdio","command":"auggie","args":["--mcp","--mcp-auto-workspace"],"env":{"AUGMENT_API_TOKEN":"'"${AUGMENT_API_TOKEN:-}"'","AUGMENT_API_URL":"'"${AUGMENT_API_URL:-}"'"}}'
echo "  auggie: ok"

# yunxiao
remove_if_exists yunxiao
claude mcp add-json --scope user yunxiao \
  '{"type":"stdio","command":"yunxiao-mcp","args":[],"env":{"YUNXIAO_ACCESS_TOKEN":"'"${YUNXIAO_ACCESS_TOKEN:-}"'"}}'
echo "  yunxiao: ok"

# exa (HTTP)
remove_if_exists exa
claude mcp add-json --scope user exa '{"type":"http","url":"https://mcp.exa.ai/mcp"}'
echo "  exa: ok"

echo ""
echo "Done. Run 'claude mcp list' to verify."
