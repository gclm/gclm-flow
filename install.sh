#!/usr/bin/env bash
# -*- mode: sh; sh-basic-offset: 4; indent-tabs-mode: nil; coding: utf-8 -*-
# vim: set filetype=sh sw=4 sts=4 et:

# gclm-flow 全局配置安装脚本
# 将自定义 agents 安装到 ~/.claude/agents/，使其在所有项目中可用
# 同时下载并安装 gclm-engine 二进制文件

set -euo pipefail

# 颜色定义（使用 printf 确保跨 shell 兼容）
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# ============================================================================
# gclm-engine 二进制安装配置
# ============================================================================

ENGINE_REPO="gclm/gclm-flow"
ENGINE_INSTALL_DIR="$HOME/.gclm-flow"
ENGINE_BINARY_NAME="gclm-engine"
ENGINE_VERSION="${GCLM_VERSION:-latest}"

# ============================================================================
# 辅助函数
# ============================================================================

# 辅助函数：彩色输出（兼容 zsh 和 bash）
print_color() {
    local color="$1"
    shift
    printf "${color}%s${NC}\n" "$*"
}

# ============================================================================
# gclm-engine 安装函数
# ============================================================================

# 检测平台（用于下载二进制）
detect_engine_platform() {
    local os=""
    local arch=""

    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            os="windows"
            ;;
        *)
            echo "Unknown"
            return
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            echo "Unknown"
            return
            ;;
    esac

    echo "${os}-${arch}"
}

# 获取 gclm-engine 最新版本
get_engine_version() {
    if [ "$ENGINE_VERSION" = "latest" ]; then
        # 使用 GitHub API 获取最新版本
        local tag
        tag=$(curl -s "https://api.github.com/repos/${ENGINE_REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$tag" ]; then
            # 如果 API 失败，使用本地版本
            echo "latest"
            return
        fi
        echo "$tag"
    else
        echo "$ENGINE_VERSION"
    fi
}

# 安装 gclm-engine 二进制
install_gclm_engine() {
    print_color "$BLUE" "=== 安装 gclm-engine ==="

    # 检测平台
    local platform
    platform=$(detect_engine_platform)

    if [ "$platform" = "Unknown" ]; then
        print_color "$YELLOW" "⚠ 跳过 gclm-engine 安装（不支持的平台）"
        return
    fi

    print_color "$BLUE" "检测到平台: $platform"

    # 创建安装目录
    mkdir -p "$ENGINE_INSTALL_DIR"

    # 检查是否存在本地编译的二进制
    if [ -f "$PROJECT_DIR/gclm-engine/build/gclm-engine" ]; then
        print_color "$BLUE" "使用本地编译的二进制..."
        cp "$PROJECT_DIR/gclm-engine/build/gclm-engine" "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME"
        chmod +x "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME"
        print_color "$GREEN" "✓ 二进制已安装"
    else
        # 从 GitHub Releases 下载
        local version
        version=$(get_engine_version)
        print_color "$BLUE" "版本: $version"

        if [ "$version" = "latest" ]; then
            print_color "$YELLOW" "⚠ 无法获取版本信息，使用本地构建"
            print_color "$YELLOW" "  请运行: cd gclm-engine && go build -o ../$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME"
            return
        fi

        local download_url="https://github.com/${ENGINE_REPO}/releases/download/${version}/${ENGINE_BINARY_NAME}-${platform}"

        print_color "$BLUE" "下载二进制..."
        if ! curl -fsSL "$download_url" -o "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME" 2>/dev/null; then
            print_color "$YELLOW" "⚠ 下载失败，跳过二进制安装"
            print_color "$YELLOW" "  手动下载: $download_url"
            return
        fi

        chmod +x "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME"
        print_color "$GREEN" "✓ 二进制已下载并安装"
    fi

    # 同步工作流文件
    sync_engine_workflows

    # 创建示例工作流
    create_example_workflow

    # 初始化数据库
    init_engine_database

    print_color "$GREEN" "✓ gclm-engine 安装完成"
}

# 同步工作流文件
sync_engine_workflows() {
    print_color "$BLUE" "同步工作流文件..."

    local source_dir=""

    # 检测工作流目录位置
    if [ -d "$PROJECT_DIR/gclm-engine/workflows" ]; then
        source_dir="$PROJECT_DIR/gclm-engine/workflows"
    elif [ -d "$PROJECT_DIR/workflows" ]; then
        source_dir="$PROJECT_DIR/workflows"
    else
        print_color "$YELLOW" "⚠ 未找到工作流目录"
        return
    fi

    # 创建目标目录
    mkdir -p "$ENGINE_INSTALL_DIR/workflows"

    # 复制工作流文件
    local count=0
    for yaml in "$source_dir"/*.yaml; do
        if [ -f "$yaml" ]; then
            cp "$yaml" "$ENGINE_INSTALL_DIR/workflows/"
            count=$((count + 1))
        fi
    done

    if [ $count -gt 0 ]; then
        print_color "$GREEN" "✓ 已同步 $count 个工作流文件"
    fi
}

# 创建示例工作流
create_example_workflow() {
    print_color "$BLUE" "创建示例工作流..."

    local example_dir="$ENGINE_INSTALL_DIR/workflows/examples"
    mkdir -p "$example_dir"

    # 创建简单自定义工作流示例
    if [ ! -f "$example_dir/custom_simple.yaml" ]; then
        cat > "$example_dir/custom_simple.yaml" << 'EOF'
# 自定义简单工作流示例
name: custom_simple
display_name: "自定义简单工作流"
description: "一个最小化的自定义工作流示例"
version: "1.0"
workflow_type: CODE_SIMPLE

nodes:
  - ref: discovery
    display_name: 需求发现
    agent: investigator
    model: haiku
    required: true
    timeout: 60

  - ref: clarification
    display_name: 澄清确认
    agent: investigator
    model: haiku
    required: true
    timeout: 60
    depends_on:
      - discovery

  - ref: implementation
    display_name: 实现
    agent: worker
    model: sonnet
    required: true
    timeout: 300
    depends_on:
      - clarification

  - ref: summary
    display_name: 总结
    agent: investigator
    model: haiku
    required: true
    timeout: 60
    depends_on:
      - implementation
EOF
        print_color "$GREEN" "✓ 已创建示例工作流: custom_simple.yaml"
    fi
}

# 初始化 gclm-engine 数据库
init_engine_database() {
    print_color "$BLUE" "初始化数据库..."

    # 运行一次命令来触发数据库初始化
    if "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME" workflow list &>/dev/null; then
        print_color "$GREEN" "✓ 数据库初始化成功"
    else
        print_color "$YELLOW" "⚠ 数据库将在首次运行时自动初始化"
    fi
}

# 更新 PATH 配置
update_engine_path() {
    local shell_config=""
    case "$SHELL" in
        */zsh)
            shell_config="$HOME/.zshrc"
            ;;
        */bash)
            shell_config="$HOME/.bashrc"
            ;;
        *)
            shell_config="$HOME/.profile"
            ;;
    esac

    # 检查是否已在 PATH 中
    if echo ":$PATH:" | grep -q ":${ENGINE_INSTALL_DIR}:"; then
        return
    fi

    # 添加到 PATH
    if [ -w "$shell_config" ]; then
        echo "" >> "$shell_config"
        echo "# gclm-engine" >> "$shell_config"
        echo "export PATH=\"\$PATH:${ENGINE_INSTALL_DIR}\"" >> "$shell_config"
        print_color "$GREEN" "✓ 已添加到 PATH: $shell_config"
        print_color "$YELLOW" "  请运行 'source $shell_config' 或重新打开终端"
    fi
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

# 全局配置模式：项目目录
PROJECT_DIR="$SCRIPT_DIR"
OS="$(detect_os)"
# 备份现有配置（只保留最新一个）
CLAUDE_DIR="$HOME/.claude"
BACKUP_DIR="$HOME/.claude.backup"

echo "🚀 安装 gclm-flow 全局配置..."
echo ""
print_color "$BLUE" "项目目录: $PROJECT_DIR"
print_color "$BLUE" "操作系统: $OS"
echo ""

# ============================================================================
# 第一步：安装 gclm-engine 二进制
# ============================================================================
install_gclm_engine
update_engine_path
echo ""

# ============================================================================
# 第二步：安装全局配置
# ============================================================================

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
cp -R "$PROJECT_DIR/agents/"*.md "$CLAUDE_DIR/agents/" 2>/dev/null || true

# 复制 commands
print_color "$BLUE" "安装 commands..."
cp -R "$PROJECT_DIR/commands/"*.md "$CLAUDE_DIR/commands/" 2>/dev/null || true

# 复制 skills
print_color "$BLUE" "安装 skills..."
cp -R "$PROJECT_DIR/skills/"* "$CLAUDE_DIR/skills/" 2>/dev/null || true

# 复制 rules
print_color "$BLUE" "安装 rules..."
cp -R "$PROJECT_DIR/rules/"*.md "$CLAUDE_DIR/rules/" 2>/dev/null || true

# 复制 hooks
print_color "$BLUE" "安装 hooks..."
cp -R "$PROJECT_DIR/hooks/"*.sh "$CLAUDE_DIR/hooks/" 2>/dev/null || true
# 设置可执行权限
chmod +x "$CLAUDE_DIR/hooks/"*.sh 2>/dev/null || true

# 检查 auggie（推荐依赖）
print_color "$BLUE" "检查 auggie..."
if command -v auggie &>/dev/null && auggie --version &>/dev/null; then
    print_color "$GREEN" "✓ auggie 已安装"
else
    print_color "$YELLOW" "⚠ auggie 未安装（推荐）"
    echo "  代码搜索增强，安装命令:"
    echo "    npm install -g @augmentcode/auggie@prerelease"
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
SETTINGS_EXAMPLE="$PROJECT_DIR/settings.example.json"
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
if [ -f "$PROJECT_DIR/CLAUDE.example.md" ]; then
    print_color "$BLUE" "安装 CLAUDE.md..."
    cp "$PROJECT_DIR/CLAUDE.example.md" "$CLAUDE_DIR/CLAUDE.md"
fi

echo ""
print_color "$GREEN" "✅ gclm-flow 全局配置安装完成！"
echo ""
echo "已安装的组件："
echo "  - gclm-engine: $([ -x "$ENGINE_INSTALL_DIR/$ENGINE_BINARY_NAME" ] && echo "✓ 已安装" || echo "✗ 未安装")"
echo "  - Agents: $(ls "$CLAUDE_DIR/agents" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Commands: $(ls "$CLAUDE_DIR/commands" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Skills: $(find "$CLAUDE_DIR/skills" -name "SKILL.md" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Rules: $(ls "$CLAUDE_DIR/rules" 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - Hooks: $(ls "$CLAUDE_DIR/hooks" 2>/dev/null | wc -l | tr -d ' ') 个"
echo ""
echo "可用的 Sub-Agents (通过自然语言调用):"
echo "  - investigator: 代码库调查（Haiku）"
echo "  - architect: 架构设计（Opus）"
echo "  - worker: 任务执行（Sonnet）"
echo "  - tdd-guide: TDD 指导（Sonnet）"
echo "  - spec-guide: SpecDD 指导（Opus）"
echo "  - code-reviewer: 代码审查（Sonnet）"
echo ""
echo "使用方法："
echo "  gclm-engine 工作流:"
echo "    gclm-engine workflow list"
echo "    gclm-engine workflow start \"修复登录页面 bug\""
echo "    gclm-engine task current <task-id>"
echo "    gclm-engine task complete <task-id> <phase-id> --output \"...\""
echo ""
echo "  gclm-flow Skills:"
echo "    /gclm <任务>"
echo "    /investigate <问题>"
echo "    /tdd <功能>"
echo "    /spec <功能>"
echo ""
if [ -n "${BACKUP_DIR:-}" ]; then
    echo "备份位置: $BACKUP_DIR"
fi
