---
name: exa-search
description: Use when you need web search, code examples, API docs, company research, or people lookup via the Exa MCP. Trigger on: search for, look up, find, what's the latest on, 搜索, 查一下, 最新进展.
---

# Exa 搜索

通过 Exa MCP 进行网页、代码、公司、人物等神经搜索。

## MCP 要求

exa MCP 已配置（`claude mcp list` 可见 exa）。

## 核心工具

### web_search_exa
通用网页搜索，适合当前信息、新闻、事实查询。

```
web_search_exa(query: "latest AI developments 2026", numResults: 5)
```

### web_search_advanced_exa
带域名和时间过滤的精确搜索。

```
web_search_advanced_exa(
  query: "React Server Components best practices",
  numResults: 5,
  includeDomains: ["github.com", "react.dev"],
  startPublishedDate: "2025-01-01"
)
```

### get_code_context_exa
从 GitHub、Stack Overflow、文档站搜索代码示例。

```
get_code_context_exa(query: "Python asyncio patterns", tokensNum: 3000)
```

tokensNum 建议：聚焦片段用 1000-2000，全面上下文用 5000+。

### company_research_exa
公司情报和新闻调研。

```
company_research_exa(companyName: "Anthropic", numResults: 5)
```

### crawling_exa
抓取指定 URL 的完整页面内容。

```
crawling_exa(url: "https://example.com/article", tokensNum: 5000)
```

### deep_researcher_start / deep_researcher_check
启动异步 AI 调研 agent，适合需要综合分析的复杂主题。

```
# 启动
deep_researcher_start(query: "comprehensive analysis of AI code editors in 2026")

# 查询结果
deep_researcher_check(researchId: "<id>")
```

## 使用模式

### 快速查询
```
web_search_exa(query: "Node.js 22 new features", numResults: 3)
```

### 代码调研
```
get_code_context_exa(query: "Rust error handling Result type", tokensNum: 3000)
```

### 公司尽调
```
company_research_exa(companyName: "Vercel", numResults: 5)
web_search_advanced_exa(query: "Vercel funding 2026", numResults: 3)
```

### 技术深研
```
deep_researcher_start(query: "WebAssembly component model adoption 2026")
# 执行其他工作...
deep_researcher_check(researchId: "<id>")
```

## 使用原则

- 优先本地代码搜索，exa 只用于需要外部信息或新鲜事实的场景
- 使用 exa 时在回复中说明使用了哪个工具及原因
- 搜索结果作为证据，不替代本地代码阅读

## 联动技能

- `deep-research`
- `brainstorming`
