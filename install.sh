#!/usr/bin/env bash
# -*- mode: sh; sh-basic-offset: 4; indent-tabs-mode: nil; coding: utf-8 -*-
# vim: set filetype=sh sw=4 sts=4 et:

# gclm-flow 插件安装脚本
# 兼容 bash 和 zsh

set -euo pipefail

# 颜色定义（使用 printf 确保跨 shell 兼容）
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 辅助函数：彩色输出（兼容 zsh 和 bash）
print_color() {
    local color="$1"
    shift
    printf "${color}%s${NC}\n" "$*"
}

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin) echo "macOS" ;;
        Linux) echo "Linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "Windows" ;;
        *) echo "Unknown" ;;
    esac
}

# 获取脚本所在目录（兼容 bash 和 zsh）
if [ -n "${BASH_SOURCE:-}" ]; then
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
else
    # zsh 兼容
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    # 如果是软链接，解析真实路径
    [ -L "$0" ] && SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)" 2>/dev/null || true
fi
# install.sh 位于项目根目录，所以 PLUGIN_DIR 就是 SCRIPT_DIR
PLUGIN_DIR="$SCRIPT_DIR"
OS="$(detect_os)"

echo "🚀 安装 gclm-flow 融合开发工作流插件..."
echo ""
print_color "$BLUE" "插件目录: $PLUGIN_DIR"
print_color "$BLUE" "操作系统: $OS"
echo ""

# 备份现有配置（只保留最新一个）
CLAUDE_DIR="$HOME/.claude"
BACKUP_DIR="$HOME/.claude.backup"

if [ -d "$CLAUDE_DIR" ]; then
    print_color "$YELLOW" "备份现有配置到: $BACKUP_DIR"
    rm -rf "$BACKUP_DIR"
    cp -R "$CLAUDE_DIR" "$BACKUP_DIR"
fi

# 创建目录结构
print_color "$BLUE" "创建目录结构..."
mkdir -p "$CLAUDE_DIR"/{agents,commands,skills,rules,hooks}

# 清空需要更新的目录（避免旧文件残留）
print_color "$BLUE" "清空旧文件..."
for dir in agents commands skills rules hooks scripts; do
    rm -rf "$CLAUDE_DIR/$dir"/*
    mkdir -p "$CLAUDE_DIR/$dir"
done

# 复制 agents
print_color "$BLUE" "安装 agents..."
cp -R "$PLUGIN_DIR/agents/"*.md "$CLAUDE_DIR/agents/" 2>/dev/null || true

# 复制 commands
print_color "$BLUE" "安装 commands..."
cp -R "$PLUGIN_DIR/commands/"*.md "$CLAUDE_DIR/commands/" 2>/dev/null || true

# 复制 skills
print_color "$BLUE" "安装 skills..."
cp -R "$PLUGIN_DIR/skills/"* "$CLAUDE_DIR/skills/" 2>/dev/null || true

# 复制 rules
print_color "$BLUE" "安装 rules..."
cp -R "$PLUGIN_DIR/rules/"*.md "$CLAUDE_DIR/rules/" 2>/dev/null || true

# 复制 hooks
print_color "$BLUE" "安装 hooks..."
mkdir -p "$CLAUDE_DIR/hooks/setup"
cp -R "$PLUGIN_DIR/hooks/"*.sh "$CLAUDE_DIR/hooks/" 2>/dev/null || true
# Setup hooks 需要放到 setup/ 目录
cp -R "$PLUGIN_DIR/hooks/setup-"*.sh "$CLAUDE_DIR/hooks/setup/" 2>/dev/null || true
chmod +x "$CLAUDE_DIR/hooks/"*.sh 2>/dev/null || true
chmod +x "$CLAUDE_DIR/hooks/setup/"*.sh 2>/dev/null || true

# 检查 auggie（必需依赖）
print_color "$BLUE" "检查 auggie..."
if command -v auggie &>/dev/null && auggie --version &>/dev/null; then
    print_color "$GREEN" "✓ auggie 已安装"
else
    print_color "$RED" "✗ auggie 未安装（必需依赖）"
    echo "  请运行: npm install -g @augmentcode/auggie@prerelease"
fi

# 检查 Playwright（可选，浏览器自动化）
print_color "$BLUE" "检查 Playwright..."
if npx @playwright/mcp --help &>/dev/null; then
    print_color "$GREEN" "✓ Playwright MCP 已安装"
else
    print_color "$YELLOW" "⚠ Playwright MCP 未安装（可选）"
    echo "  用于跨浏览器测试，安装命令:"
    echo "    npm install -g @playwright/mcp"
    echo "    npx playwright install"
fi

# 检查 Chrome DevTools MCP（推荐，Chrome 调试）
print_color "$BLUE" "检查 Chrome DevTools MCP..."
if npx chrome-devtools-mcp --help &>/dev/null; then
    print_color "$GREEN" "✓ Chrome DevTools MCP 已安装"
else
    print_color "$YELLOW" "⚠ Chrome DevTools MCP 未安装（推荐）"
    echo "  Google 官方 Chrome 调试工具，安装命令:"
    echo "    npm install -g chrome-devtools-mcp"
fi

# 配置 MCP 服务器和 hooks 到全局 settings.json
print_color "$BLUE" "配置 settings.json..."
SETTINGS_EXAMPLE="$PLUGIN_DIR/settings.example.json"
SETTINGS_FILE="$CLAUDE_DIR/settings.json"

if [ -f "$SETTINGS_EXAMPLE" ]; then
    # 备份现有 settings.json
    if [ -f "$SETTINGS_FILE" ]; then
        cp "$SETTINGS_FILE" "$HOME/.claude-settings-backup.json"
    fi

    # 使用 Python 合并配置
    if command -v python3 &>/dev/null; then
        python3 - "$SETTINGS_FILE" "$SETTINGS_EXAMPLE" <<'PYTHON_SCRIPT'
import json
import sys

settings_file = sys.argv[1]
example_file = sys.argv[2]

# 读取现有 settings.json
try:
    with open(settings_file, 'r') as f:
        settings = json.load(f)
except FileNotFoundError:
    settings = {}

# 读取示例配置
with open(example_file, 'r') as f:
    example = json.load(f)

# 合并 mcpServers
if 'mcpServers' not in settings:
    settings['mcpServers'] = {}
settings['mcpServers'].update(example.get('mcpServers', {}))

# 合并 hooks
if 'hooks' not in settings:
    settings['hooks'] = {}
settings['hooks'].update(example.get('hooks', {}))

# 写回 settings.json
with open(settings_file, 'w') as f:
    json.dump(settings, f, indent=2)

# 输出已配置的服务器
servers = list(example.get('mcpServers', {}).keys())
for server in servers:
    print(f"    - {server}")
PYTHON_SCRIPT

        print_color "$GREEN" "✓ settings.json 已更新"
        echo ""
        echo "  已配置的 MCP 服务器:"
    else
        print_color "$YELLOW" "⚠ 无法自动配置（需要 Python 3）"
        echo "  请手动复制 $SETTINGS_EXAMPLE 到 ~/.claude/settings.json"
    fi
else
    print_color "$YELLOW" "⚠ settings.example.json 不存在"
fi

# macOS 专属：配置通知功能
if [ "$OS" = "macOS" ]; then
    print_color "$BLUE" "配置通知功能 (macOS)..."

    if command -v terminal-notifier &>/dev/null; then
        print_color "$GREEN" "✓ terminal-notifier 已安装"
    else
        print_color "$YELLOW" "⚠ terminal-notifier 未安装（可选）"
        echo "  请运行: brew install terminal-notifier"
    fi

    # 检查 ClaudeNotifier.app
    if [ -d "/Applications/ClaudeNotifier.app" ]; then
        print_color "$GREEN" "✓ ClaudeNotifier.app 已安装"
    else
        print_color "$YELLOW" "⚠ ClaudeNotifier.app 未安装（可选）"
        echo "  请将 ClaudeNotifier.app 放到 /Applications/ 目录"
    fi
else
    print_color "$BLUE" "跳过通知功能配置（仅支持 macOS）"
fi

# 复制 CLAUDE.example.md 到 CLAUDE.md
if [ -f "$PLUGIN_DIR/CLAUDE.example.md" ]; then
    print_color "$BLUE" "安装 CLAUDE.md..."
    cp "$PLUGIN_DIR/CLAUDE.example.md" "$CLAUDE_DIR/CLAUDE.md"
fi

echo ""
print_color "$GREEN" "✅ gclm-flow 安装完成！"
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
echo "  /tdd <功能>          - 测试驱动开发 (TDD)"
echo "  /spec <功能>         - 规范驱动开发 (SpecDD)"
echo "  /llmdoc              - 更新项目文档"
echo ""
if [ -n "${BACKUP_DIR:-}" ]; then
    echo "备份位置: $BACKUP_DIR"
fi
