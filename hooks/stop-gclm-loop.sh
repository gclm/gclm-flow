#!/bin/bash
# gclm-flow Stop Hook
# 当工作流未完成时阻止中途退出

STATE_FILES=$(find .claude -name "gclm.*.local.md" 2>/dev/null || true)

# 检查是否有活跃的状态文件
if [ -z "$STATE_FILES" ]; then
    exit 0
fi

# 检查每个状态文件
for state_file in $STATE_FILES; do
    # 读取状态
    ACTIVE=$(grep "^active:" "$state_file" 2>/dev/null | cut -d':' -f2 | xargs || echo "false")
    CURRENT_PHASE=$(grep "^current_phase:" "$state_file" 2>/dev/null | cut -d':' -f2 | xargs || echo "0")
    MAX_PHASES=$(grep "^max_phases:" "$state_file" 2>/dev/null | cut -d':' -f2 | xargs || echo "8")

    # 如果未激活或已完成，跳过
    if [ "$ACTIVE" != "true" ] || [ "$CURRENT_PHASE" = "$MAX_PHASES" ]; then
        continue
    fi

    # 工作流进行中，警告用户
    echo ""
    echo "=========================================="
    echo "  gclm-flow 工作流正在进行中"
    echo "=========================================="
    echo ""
    echo "状态文件: $state_file"
    echo "当前阶段: Phase $CURRENT_PHASE / $MAX_PHASES"
    echo ""
    echo "中途退出会导致："
    echo "  - 上下文丢失"
    echo "  - 状态不一致"
    echo "  - 需要手动恢复"
    echo ""
    echo "如需强制退出，请运行："
    echo "  sed -i.bak 's/^active: true/active: false/' $state_file"
    echo ""
    echo "=========================================="
    exit 1
done

exit 0
