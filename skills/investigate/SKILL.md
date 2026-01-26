---
name: investigate
description: "快速代码库调查 - 文档优先，无状态输出。用于'什么是'、'X怎么工作'、'分析'等问题"
disable-model-invocation: false
context: fork
allowed-tools: Read, Glob, Grep, Bash, WebSearch, WebFetch
---

# /investigate Skill

快速、文档优先的代码库调查。

## 触发条件

当用户询问：
- "什么是..."
- "X是怎么工作的"
- "分析一下..."
- "解释一下..."

## 调查协议

### Phase 1: Documentation First (强制)

**在接触任何源代码前，必须：**

1. 检查 `llmdoc/` 是否存在
2. 如果存在，按顺序读取：
   - `llmdoc/index.md` - 导航和概览
   - `llmdoc/overview/*.md` - 项目上下文
   - `llmdoc/architecture/*.md` - 系统设计
   - `llmdoc/guides/*.md` - 工作流程
   - `llmdoc/reference/*.md` - 约定和规范

### Phase 2: Code Investigation

**仅在文档不足时，调查源代码：**

1. 使用 `Glob` 查找相关文件
2. 使用 `Grep` 搜索模式
3. 使用 `Read` 检查特定文件

### Phase 3: Report

输出简洁报告，结构如下：

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

## 关键实践

- **Stateless**: 直接输出，不写入文件
- **Concise**: 报告 < 150 行
- **No Code Blocks**: 使用 `path/file.ext` 格式引用，不粘贴代码
- **Objective**: 只陈述事实，无主观判断

## 输出示例

```
#### Code Sections
- `src/auth/login.ts:15~45` (handleLogin): 处理用户登录
- `src/middleware/auth.ts:8~30` (authMiddleware): JWT 验证中间件
- `src/db/models/user.ts:12~50` (User): 用户数据模型

#### Report

**Conclusions:**
> 认证系统使用 JWT tokens，登录时生成包含 userId 和 role 的 token。

**Relations:**
> login.ts → authMiddleware (token 验证)
> login.ts → User 模型 (用户查询)

**Result:**
> 用户登录流程：提交凭证 → 验证 → 生成 JWT → 返回 token。后续请求通过 authMiddleware 验证 token。
```
