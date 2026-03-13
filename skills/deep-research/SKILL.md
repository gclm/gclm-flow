---
name: deep-research
description: Use when the user needs thorough, cited research on any topic from multiple web sources. Trigger on: research, deep dive, investigate, what's the current state of, 调研, 深度研究.
---

# 深度调研

多源网络调研，输出带引用的结构化报告。使用 exa MCP 作为主要搜索工具。

## 工作流程

### 1. 明确目标

如有必要，快速确认：
- 目标是学习、决策还是输出内容？
- 是否有特定角度或深度要求？

用户说"直接研究"时跳过，用合理默认值继续。

### 2. 拆分子问题

将主题拆成 3-5 个研究子问题，每个方向独立可搜索。

### 3. 多源搜索

每个子问题使用 exa MCP 搜索：

```
# 通用搜索
web_search_exa(query: "<关键词>", numResults: 8)

# 限定来源或时间
web_search_advanced_exa(
  query: "<关键词>",
  numResults: 5,
  includeDomains: ["github.com", "docs.xxx.com"],
  startPublishedDate: "2025-01-01"
)

# 代码示例搜索
get_code_context_exa(query: "<技术关键词>", tokensNum: 3000)
```

目标：15-30 个不重复来源，优先学术、官方、权威媒体。

### 4. 精读关键来源

对最有价值的 3-5 个 URL 抓取全文：

```
crawling_exa(url: "<url>", tokensNum: 5000)
```

不能仅依赖搜索摘要。

### 5. 输出报告

```markdown
# [主题]：调研报告
*时间: [日期] | 来源数: N | 置信度: 高/中/低*

## 核心摘要
[3-5 句关键发现]

## 1. [主题一]
- 关键发现（[来源名称](url)）
- 数据佐证（[来源名称](url)）

## 2. [主题二]
...

## 关键结论
- 可执行洞察 1
- 可执行洞察 2

## 来源列表
1. [标题](url) — 一句话摘要
2. ...

## 调研方法
搜索了 N 个查询，分析了 M 个来源。
子问题：[列表]
```

## 并行调研

主题较宽时，用 Agent tool 并行启动多个调研子 agent：
- Agent 1：子问题 1-2
- Agent 2：子问题 3-4
- Agent 3：子问题 5 + 交叉主题

主 session 汇总结果输出最终报告。

## 质量规则

1. 每个结论必须有来源，不允许无引用断言。
2. 单一来源的说法需标注"待交叉验证"。
3. 优先最近 12 个月的来源。
4. 信息不足时明确说明，不猜测。
5. 区分事实、推断和观点。

## 联动技能

- `exa-search`
- `brainstorming`
