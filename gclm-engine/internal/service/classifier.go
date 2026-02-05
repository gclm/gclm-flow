package service

import "strings"

// WorkflowClassifier 工作流类型分类器
type WorkflowClassifier struct {
	// 关键词配置
	analyzePhrases   []string
	analyzeKeywords  []string
	docPhrases       []string
	docKeywords      []string
	bugPhrases       []string
	bugKeywords      []string
	featureKeywords  []string
}

// NewWorkflowClassifier 创建默认分类器
func NewWorkflowClassifier() *WorkflowClassifier {
	return &WorkflowClassifier{
		// 分析类 - 优先级最高
		analyzePhrases: []string{
			"分析代码", "代码分析", "问题分析", "分析问题",
			"诊断", "代码诊断", "性能分析", "性能评估",
			"安全审计", "安全分析", "代码审查", "代码review",
			"评估代码", "检查代码", "分析一下", "帮我分析",
			"看看代码", "理解代码", "分析架构", "架构分析",
		},
		analyzeKeywords: []string{
			"分析", "诊断", "审计", "评估", "检查", "review",
		},
		docPhrases: []string{
			"编写文档", "文档编写", "方案设计", "设计文档",
			"需求分析", "技术方案", "架构设计", "api文档", "spec文档",
			"写文档", "写方案", "写设计", "文档化",
		},
		docKeywords: []string{"编写文档", "文档编写", "方案", "设计文档", "spec文档"},
		bugPhrases: []string{
			"修复bug", "fix bug", "bug修复", "修复错误", "解决bug",
		},
		bugKeywords: []string{"bug", "修复", "fix error", "error fix", "调试", "debug"},
		featureKeywords: []string{"功能", "模块", "开发", "重构", "实现"},
	}
}

// Classify 检测工作流类型
// 返回: "ANALYZE", "DOCUMENT", "CODE_SIMPLE", "CODE_COMPLEX"
func (wc *WorkflowClassifier) Classify(prompt string) string {
	score := 0
	lowerPrompt := strings.ToLower(prompt)

	// === 分析类检测（优先级最高） ===
	// 分析类短语（+10分，确保优先匹配）
	for _, phrase := range wc.analyzePhrases {
		if strings.Contains(lowerPrompt, strings.ToLower(phrase)) {
			score += 10
		}
	}

	// 分析类单词（+5分）
	for _, kw := range wc.analyzeKeywords {
		if contains(prompt, kw) {
			score += 5
		}
	}

	// 如果分析分数 > 0，直接返回 ANALYZE
	if score > 0 {
		// 但要排除"编写文档"这样的写文档场景
		for _, docPhrase := range wc.docPhrases {
			if strings.Contains(lowerPrompt, strings.ToLower(docPhrase)) {
				return "DOCUMENT"
			}
		}
		return "ANALYZE"
	}

	// === 其他类型检测 ===
	score = 0 // 重置分数

	// 文档类短语（+5分）
	for _, phrase := range wc.docPhrases {
		if strings.Contains(lowerPrompt, strings.ToLower(phrase)) {
			score += 5
		}
	}

	// 文档类单词（+3分）
	for _, kw := range wc.docKeywords {
		if contains(prompt, kw) {
			score += 3
		}
	}

	// Bug修复短语（-5分）
	for _, phrase := range wc.bugPhrases {
		if strings.Contains(lowerPrompt, strings.ToLower(phrase)) {
			score -= 5
		}
	}

	// Bug修复单词（-3分）
	for _, kw := range wc.bugKeywords {
		if contains(prompt, kw) {
			score -= 3
		}
	}

	// 功能开发单词（-1分）
	for _, kw := range wc.featureKeywords {
		if contains(prompt, kw) {
			score -= 1
		}
	}

	// 分类阈值
	if score >= 3 {
		return "DOCUMENT"
	} else if score <= -3 {
		return "CODE_SIMPLE"
	}
	return "CODE_COMPLEX"
}

// DetectWorkflowType 便捷函数，使用默认分类器
func DetectWorkflowType(prompt string) string {
	return NewWorkflowClassifier().Classify(prompt)
}

// contains 检查字符串是否包含中文或英文关键词
func contains(s, substr string) bool {
	return strings.Contains(s, substr) ||
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
