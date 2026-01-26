#!/bin/bash

# gclm-flow 插件安装脚本

set -euo pipefail

echo "🚀 安装 gclm-flow 融合开发工作流插件..."
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_DIR="$(dirname "$SCRIPT_DIR")"

echo "${BLUE}插件目录:${NC} $PLUGIN_DIR"
echo ""

# 备份现有配置
CLAUDE_DIR="$HOME/.claude"
BACKUP_DIR="$HOME/.claude.backup.$(date +%Y%m%d%H%M%S)"

if [ -d "$CLAUDE_DIR" ]; then
    echo "${YELLOW}备份现有配置到:${NC} $BACKUP_DIR"
    cp -R "$CLAUDE_DIR" "$BACKUP_DIR"
fi

# 创建目录结构
echo "${BLUE}创建目录结构...${NC}"
mkdir -p "$CLAUDE_DIR"/{agents,commands,skills,rules,hooks,scripts}

# 复制 agents
echo "${BLUE}安装 agents...${NC}"
cp -R "$PLUGIN_DIR/agents/"*.md "$CLAUDE_DIR/agents/" 2>/dev/null || true

# 复制 commands
echo "${BLUE}安装 commands...${NC}"
cp -R "$PLUGIN_DIR/commands/"*.md "$CLAUDE_DIR/commands/" 2>/dev/null || true

# 复制 skills
echo "${BLUE}安装 skills...${NC}"
cp -R "$PLUGIN_DIR/skills/"* "$CLAUDE_DIR/skills/" 2>/dev/null || true

# 复制 rules
echo "${BLUE}安装 rules...${NC}"
cp -R "$PLUGIN_DIR/rules/"*.md "$CLAUDE_DIR/rules/" 2>/dev/null || true

# 复制 hooks
echo "${BLUE}安装 hooks...${NC}"
cp -R "$PLUGIN_DIR/hooks/"* "$CLAUDE_DIR/hooks/" 2>/dev/null || true
chmod +x "$CLAUDE_DIR/hooks/"*.sh 2>/dev/null || true

# 复制 scripts
echo "${BLUE}安装 scripts...${NC}"
cp -R "$PLUGIN_DIR/scripts/"* "$CLAUDE_DIR/scripts/" 2>/dev/null || true
chmod +x "$CLAUDE_DIR/scripts/"*.sh 2>/dev/null || true

# 复制 CLAUDE.example.md 到 CLAUDE.md
if [ -f "$PLUGIN_DIR/CLAUDE.example.md" ]; then
    echo "${BLUE}安装 CLAUDE.md...${NC}"
    cp "$PLUGIN_DIR/CLAUDE.example.md" "$CLAUDE_DIR/CLAUDE.md"
fi

echo ""
echo "${GREEN}✅ gclm-flow 安装完成！${NC}"
echo ""
echo "已安装的组件："
echo "  - Agents: $(ls "$CLAUDE_DIR/agents" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Commands: $(ls "$CLAUDE_DIR/commands" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Skills: $(find "$CLAUDE_DIR/skills" -name "SKILL.md" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Rules: $(ls "$CLAUDE_DIR/rules" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Hooks: $(ls "$CLAUDE_DIR/hooks" 2>/dev/null | wc -l | tr -d ' ') 个"
echo ""
echo "使用方法："
echo "  /gclm <功能描述>     - 启动完整工作流"
echo "  /investigate <问题>   - 快速代码库调查"
echo "  /tdd <功能>          - 测试驱动开发"
echo ""
if [ -n "${BACKUP_DIR:-}" ]; then
    echo "备份位置: $BACKUP_DIR"
fi
