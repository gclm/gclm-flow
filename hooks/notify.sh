#!/bin/bash

# Claude Code 用户操作通知脚本
# 发送带自定义图标的 macOS 系统通知

MESSAGE="${1:-Claude Code 需要您的操作}"

# 使用 ClaudeNotifier.app 的 bundle ID
terminal-notifier \
  -message "$MESSAGE" \
  -title "Claude Code" \
  -subtitle "请检查终端" \
  -sender "com.claude.notifier" \
  -sound default

exit 0
