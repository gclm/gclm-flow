---
name: documentation
description: |
  文档管理技能。当用户要求文档、docs、README、注释、API 文档时自动触发。
  包含：(1) llmdoc 结构 (2) 文档约定 (3) API 文档生成
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - docs
    - readme
    - api
---

# 文档管理

## llmdoc 结构

```
llmdoc/
├── reference/
│   ├── conventions/      # 项目约定
│   └── guidelines/       # 开发指南
├── architecture/         # 架构文档
├── api/                  # API 文档
└── guides/               # 使用指南
```

## 文档约定

### 命名规范
- 使用小写和连字符：`api-design.md`
- 版本化：`migration-v2.md`

### 格式规范
- 使用 Markdown
- 标题层级清晰
- 代码块指定语言

### API 文档

```markdown
# API 名称

## 端点
`GET /api/v1/users`

## 参数
| 名称 | 类型 | 必需 | 描述 |
|------|------|------|------|

## 响应
```json
{
  "success": true,
  "data": {}
}
```

## 示例
```

## 详见

- [llmdoc-structure.md](references/llmdoc-structure.md)
- [doc-conventions.md](references/doc-conventions.md)
