#!/usr/bin/env bash

set -euo pipefail

usage() {
	cat <<'EOF'
Usage: setup-gclm.sh [options] PROMPT...

创建（或覆盖）项目状态文件：
  .claude/gclm.{task_id}.local.md

Options:
  --max-phases N            默认: 8
  --completion-promise STR  默认: <promise>GCLM_WORKFLOW_COMPLETE</promise>
  -h, --help                显示此帮助
EOF
}

die() {
	echo "❌ $*" >&2
	exit 1
}

phase_name_for() {
	case "${1:-}" in
		0) echo "llmdoc Reading" ;;
		1) echo "Discovery" ;;
		2) echo "Exploration" ;;
		3) echo "Clarification" ;;
		4) echo "Architecture" ;;
		5) echo "TDD Red" ;;
		6) echo "TDD Green" ;;
		7) echo "Refactor + Doc" ;;
		8) echo "Summary" ;;
		*) echo "Phase ${1:-unknown}" ;;
	esac
}

max_phases=8
completion_promise="<promise>GCLM_WORKFLOW_COMPLETE</promise>"
declare -a prompt_parts=()

while [ $# -gt 0 ]; do
	case "$1" in
		-h|--help)
			usage
			exit 0
			;;
		--max-phases)
			[ $# -ge 2 ] || die "--max-phases 需要值"
			max_phases="$2"
			shift 2
			;;
		--completion-promise)
			[ $# -ge 2 ] || die "--completion-promise 需要值"
			completion_promise="$2"
			shift 2
			;;
		--)
			shift
			while [ $# -gt 0 ]; do
				prompt_parts+=("$1")
				shift
			done
			break
			;;
		-*)
			die "未知参数: $1 (使用 --help)"
			;;
		*)
			prompt_parts+=("$1")
			shift
			;;
	esac
done

prompt="${prompt_parts[*]:-}"
[ -n "$prompt" ] || die "PROMPT 是必需的 (使用 --help)"

if ! [[ "$max_phases" =~ ^[0-9]+$ ]] || [ "$max_phases" -lt 1 ]; then
	die "--max-phases 必须是正整数"
fi

project_dir="${CLAUDE_PROJECT_DIR:-$PWD}"
state_dir="${project_dir}/.claude"

task_id="$(date +%s)-$$-$(head -c 4 /dev/urandom | od -An -tx1 | tr -d ' \n')"
state_file="${state_dir}/gclm.${task_id}.local.md"

mkdir -p "$state_dir"

phase_name="$(phase_name_for 0)"

cat > "$state_file" << EOF
---
active: true
current_phase: 0
phase_name: "$phase_name"
max_phases: $max_phases
completion_promise: "$completion_promise"
---

# gclm-flow loop state

## Prompt
$prompt

## 8 Phases
| Phase | Name | Status |
|:---|:---|:---|
| 0 | llmdoc Reading | in_progress |
| 1 | Discovery | pending |
| 2 | Exploration | pending |
| 3 | Clarification | pending |
| 4 | Architecture | pending |
| 5 | TDD Red | pending |
| 6 | TDD Green | pending |
| 7 | Refactor + Doc | pending |
| 8 | Summary | pending |

## Notes
- 按进度更新 frontmatter 的 current_phase/phase_name
- 完成时在最终输出中包含 frontmatter completion_promise

## Key Constraints
1. Phase 0: 强制读取 llmdoc
2. Phase 3: 不可跳过，必须澄清
3. Phase 5: TDD Red，测试必须先失败
4. Phase 7: 必须询问是否更新文档
EOF

echo "✅ Initialized: $state_file"
echo "task_id: $task_id"
echo "phase: 0/$max_phases ($phase_name)"
echo "completion_promise: $completion_promise"
