#!/usr/bin/env bash

set -euo pipefail

usage() {
	cat <<'EOF'
Usage: setup-gclm.sh [options] PROMPT...

创建（或覆盖）项目状态文件：
  .claude/gclm.{task_id}.local.md

Options:
  --export                 输出 JSON 元数据而非创建文件（供 LLM Write 工具使用）
  --max-phases N           默认: 9
  --completion-promise STR 默认: <promise>GCLM_WORKFLOW_COMPLETE</promise>
  --workflow TYPE          手动指定工作流类型 (DOCUMENT|CODE_SIMPLE|CODE_COMPLEX)
  -h, --help               显示此帮助

Examples:
  # 标准模式：直接创建状态文件
  setup-gclm.sh "实现用户登录功能"

  # 导出模式：输出 JSON 供 LLM 使用
  setup-gclm.sh --export "实现用户登录功能"
EOF
}

die() {
	echo "❌ $*" >&2
	exit 1
}

# 获取阶段名称（双语）- 9 阶段系统
phase_name_for() {
	local workflow="${2:-CODE_COMPLEX}"
	case "${1:-}" in
		0) echo "llmdoc Reading / 读取文档" ;;
		1) echo "Discovery / 需求发现" ;;
		2) echo "Exploration / 探索研究" ;;
		3) echo "Clarification / 澄清确认" ;;
		4) echo "Architecture / 架构设计" ;;
		5)
			case "$workflow" in
				DOCUMENT) echo "Draft / 起草文档" ;;
				*) echo "Spec / 规范文档" ;;
			esac
			;;
		6)
			case "$workflow" in
				DOCUMENT) echo "Refine / 完善内容" ;;
				*) echo "TDD Red / 编写测试" ;;
			esac
			;;
		7)
			case "$workflow" in
				DOCUMENT) echo "Review / 质量审查" ;;
				*) echo "TDD Green / 编写实现" ;;
			esac
			;;
		8) echo "Refactor + Security + Review / 重构审查" ;;
		9) echo "Summary / 完成总结" ;;
		*) echo "Phase ${1:-unknown}" ;;
	esac
}

# 自动检测工作流类型（改进版：短语匹配优先，减少误判）
detect_workflow_type() {
	local prompt="$1"
	local score=0

	# 文档/方案关键词（高优先级短语匹配）
	if echo "$prompt" | grep -iqE "编写文档|文档编写|方案设计|设计文档|需求分析|技术方案|架构设计|API文档|README|Spec文档|规范文档"; then
		score=$((score + 5))  # 短语匹配加分更多
	elif echo "$prompt" | grep -iqE "文档|方案|需求|分析|架构|规范|说明"; then
		score=$((score + 3))
	fi

	# Bug修复关键词（高优先级，负分）
	if echo "$prompt" | grep -iqE "修复bug|fix bug|bug修复|修复错误|解决bug"; then
		score=$((score - 5))
	elif echo "$prompt" | grep -iqE "bug|修复|fix.*error|error.*fix|调试|debug"; then
		score=$((score - 3))
	fi

	# 新功能开发关键词（中等优先级，轻微负分）
	if echo "$prompt" | grep -iqE "开发新功能|实现功能|添加功能|新模块|功能开发"; then
		score=$((score - 1))
	elif echo "$prompt" | grep -iqE "功能|模块|开发|重构|实现"; then
		score=$((score - 1))
	fi

	# 分类（调整阈值）
	if [ "$score" -ge 3 ]; then
		echo "DOCUMENT"     # 文档编写/方案设计
	elif [ "$score" -le -3 ]; then
		echo "CODE_SIMPLE"  # Bug修复/小修改
	else
		echo "CODE_COMPLEX" # 新功能/模块开发
	fi
}

# 根据工作流类型获取阶段列表 - 9 阶段系统
get_phases_for_workflow() {
	local workflow="$1"
	case "$workflow" in
		DOCUMENT)
			# 文档工作流：0 → 1 → 2 → 3 → 5(Draft) → 6(Refine) → 7(Review) → 8 → 9
			echo "0,1,2,3,5,6,7,8,9"
			;;
		CODE_SIMPLE)
			# 简单代码工作流：0 → 1 → 3 → 6(TDD Red) → 7(TDD Green) → 8 → 9
			echo "0,1,3,6,7,8,9"
			;;
		CODE_COMPLEX|*)
			# 复杂代码工作流：全部 9 个阶段
			echo "0,1,2,3,4,5,6,7,8,9"
			;;
	esac
}

max_phases=9
completion_promise="<promise>GCLM_WORKFLOW_COMPLETE</promise>"
workflow_type=""  # 允许手动指定
export_mode=false  # 导出模式：输出 JSON 而非创建文件
declare -a prompt_parts=()

# 检测代码搜索方法（分层回退）
detect_code_search_method() {
	if auggie --help &>/dev/null 2>&1; then
		echo "auggie"
	else
		echo "llmdoc+grep"
	fi
}

while [ $# -gt 0 ]; do
	case "$1" in
		-h|--help)
			usage
			exit 0
			;;
		--export)
			export_mode=true
			shift
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
		--workflow)
			[ $# -ge 2 ] || die "--workflow 需要值 (DOCUMENT|CODE_SIMPLE|CODE_COMPLEX)"
			workflow_type="$2"
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

# 确定工作流类型（手动指定 > 自动检测）
if [ -n "$workflow_type" ]; then
	# 验证手动指定的类型
	case "$workflow_type" in
		DOCUMENT|CODE_SIMPLE|CODE_COMPLEX)
			detected_workflow="$workflow_type"
			;;
		*)
			die "无效的工作流类型: $workflow_type (必须是 DOCUMENT|CODE_SIMPLE|CODE_COMPLEX)"
			;;
	esac
else
	# 自动检测
	detected_workflow="$(detect_workflow_type "$prompt")"
fi

# 检测代码搜索方法
code_search_method="$(detect_code_search_method)"

if ! [[ "$max_phases" =~ ^[0-9]+$ ]] || [ "$max_phases" -lt 1 ]; then
	die "--max-phases 必须是正整数"
fi

project_dir="${CLAUDE_PROJECT_DIR:-$PWD}"
state_dir="${project_dir}/.claude"

task_id="$(date +%s)-$$-$(head -c 4 /dev/urandom | od -An -tx1 | tr -d ' \n')"
state_file="${state_dir}/gclm.${task_id}.local.md"

mkdir -p "$state_dir"

# 自动配置项目级权限（简化版）
settings_file="${state_dir}/settings.local.json"

if [ ! -f "$settings_file" ]; then
	# 首次运行，创建基础权限配置
	# 注意：这会授权常见的 .claude/ 操作
	cat > "$settings_file" << 'EOF'
{
  "permissions": {
    "allow": [
      "Bash(\"mkdir -p .claude*\")",
      "Bash(\"ls -la .claude/\")"
    ]
  },
  "_comment": "gclm-flow 自动生成。首次使用时请在授权提示选择 'Yes, and always allow access to .claude/'"
}
EOF
	echo "ℹ️  首次使用，请在授权提示选择 'Yes, and always allow access to .claude/' 以避免重复授权"
fi

phase_name="$(phase_name_for 0 "$detected_workflow")"

# 预先计算工作流描述
workflow_desc=""
case "$detected_workflow" in
	DOCUMENT) workflow_desc="📝 文档编写/方案设计工作流" ;;
	CODE_SIMPLE) workflow_desc="🔧 简单代码工作流 (Bug修复/小修改)" ;;
	CODE_COMPLEX) workflow_desc="🚀 复杂代码工作流 (新功能/模块开发)" ;;
esac

# 预先计算代码搜索方法描述
code_search_desc=""
if [ "$code_search_method" = "auggie" ]; then
	code_search_desc="✅ auggie 语义搜索（推荐）"
else
	code_search_desc="⚠️  llmdoc + Grep 备选方案（auggie 未安装）

💡 安装 auggie 获得更好的语义搜索体验:
   npm install -g @augmentcode/auggie@prerelease"
fi

# 获取该工作流类型的阶段列表
phases_list="$(get_phases_for_workflow "$detected_workflow")"
IFS=',' read -ra phases <<< "$phases_list"

# 生成阶段表格
phases_table=""
for p in "${phases[@]}"; do
	p_name="$(phase_name_for "$p" "$detected_workflow")"
	if [ "$p" = "0" ]; then
		status="in_progress"
	else
		status="pending"
	fi
	phases_table="$phases_table| $p | $p_name | $status |
"
done

# 生成状态文件完整内容
state_file_content=$(cat <<EOF
---
active: true
current_phase: 0
phase_name: "$phase_name"
max_phases: $max_phases
workflow_type: "$detected_workflow"
code_search: "$code_search_method"
completion_promise: "$completion_promise"
---

# gclm-loop state

## Prompt
$prompt

## Workflow Type: $detected_workflow
$workflow_desc

## Code Search Method
$code_search_desc

## Phases for $detected_workflow
$phases_table
## Workflow Phases Summary (双语)
| Phase | Name / 名称 | DOCUMENT | CODE_SIMPLE | CODE_COMPLEX |
|:---:|:---|:---:|:---:|:---:|
| 0 | llmdoc Reading / 读取文档 | ✅ | ✅ | ✅ |
| 1 | Discovery / 需求发现 | ✅ | ✅ | ✅ |
| 2 | Exploration / 探索研究 | ✅ | ❌ | ✅ |
| 3 | Clarification / 澄清确认 | ✅ | ✅ | ✅ |
| 4 | Architecture / 架构设计 | ❌ | ❌ | ✅ |
| 5 | Spec/Draft / 规范/起草 | ✅ (起草) | ❌ | ✅ (Spec) |
| 6 | Refactor/TDD Red / 完善/测试 | ✅ (完善) | ✅ (TDD Red) | ✅ (TDD Red) |
| 7 | Review/TDD Green / 审查/实现 | ✅ (审查) | ✅ (TDD Green) | ✅ (TDD Green) |
| 8 | Refactor+Review / 重构审查 | ✅ | ✅ | ✅ |
| 9 | Summary / 完成总结 | ✅ | ✅ | ✅ |

## Notes
- 按进度更新 frontmatter 的 current_phase/phase_name
- 完成时在最终输出中包含 frontmatter completion_promise
- 状态更新自动进行，无需用户确认
- llmdoc 不存在时自动生成
- **Phase 3 澄清时，可手动调整 workflow_type**

## Key Constraints
1. Phase 0: 强制读取 llmdoc，代码搜索采用分层回退（auggie → llmdoc → Grep）
2. Phase 3: 不可跳过，必须澄清 + 确认/调整工作流类型
3. DOCUMENT 工作流: Phase 5 起草文档, Phase 6 完善内容, Phase 7 质量审查
4. CODE_SIMPLE: 跳过 Phase 2, 4, 5（条件：estimated_files > 3 时不跳过 Phase 2）
5. CODE_COMPLEX: 全部 9 阶段 + SpecDD
6. 状态更新: 自动化，无需授权
EOF
)

# 导出模式：输出 JSON 供 LLM 使用
if [ "$export_mode" = true ]; then
	# 输出 JSON 格式的元数据
	cat <<EOF
{
  "task_id": "$task_id",
  "state_file": "$state_file",
  "workflow_type": "$detected_workflow",
  "workflow_desc": "$workflow_desc",
  "code_search_method": "$code_search_method",
  "code_search_desc": $(echo "$code_search_desc" | jq -Rs . 2>/dev/null || echo '"(see standard mode output)"'),
  "phase_name": "$phase_name",
  "current_phase": 0,
  "max_phases": $max_phases,
  "completion_promise": "$completion_promise",
  "prompt": $(echo "$prompt" | jq -Rs . 2>/dev/null || echo '"(use escape if needed)"'),
  "state_file_content": $(echo "$state_file_content" | jq -Rs . 2>/dev/null || echo '"(too large for JSON)"')
}
EOF
	exit 0
fi

cat > "$state_file" <<< "$state_file_content"

echo "✅ Initialized: $state_file"
echo "task_id: $task_id"
echo "phase: 0/$max_phases ($phase_name)"
echo "completion_promise: $completion_promise"
