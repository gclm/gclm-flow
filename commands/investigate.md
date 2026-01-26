---
description: "快速代码库调查 - 文档优先，无状态输出"
---

# /investigate 命令

快速、文档优先的代码库调查。

## 使用方法

```
/investigate <问题>
```

## 示例

```
/investigate 认证系统是怎么工作的？
/investigate 数据库连接在哪里配置？
/investigate 解释一下用户注册流程
```

## 调查协议

### Phase 1: Documentation First (强制)

在接触任何源代码前，必须：

1. 检查 `llmdoc/` 是否存在
2. 如果存在，按顺序读取：
   - `llmdoc/index.md`
   - `llmdoc/overview/*.md`
   - `llmdoc/architecture/*.md`
   - `llmdoc/guides/*.md`
   - `llmdoc/reference/*.md`

### Phase 2: Code Investigation

仅在文档不足时，调查源代码。

### Phase 3: Report

输出简洁报告，包含：
- Code Sections (文件引用)
- Conclusions (关键发现)
- Relations (模块关系)
- Result (直接回答)

## 输出格式

```markdown
#### Code Sections
- `path/to/file.ext:line~line` (SymbolName): Brief description

#### Report

**Conclusions:**
> Key findings...

**Relations:**
> File/module relationships...

**Result:**
> Direct answer to the question...
```

## 特点

- **文档优先**: 优先读取 llmdoc
- **无状态**: 不写入文件
- **简洁**: 报告 < 150 行
- **客观**: 只陈述事实
