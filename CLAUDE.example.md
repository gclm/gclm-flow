# Claude Code 全局配置

本文件为 Claude Code 提供全局通用的工作指导，适用于所有项目。

---

## 核心开发原则

**SpecDD + TDD + Document-First + 智能分流 + 多 Agent 并行**

- **SpecDD**: 复杂模块先写规范文档
- **TDD**: 测试驱动开发 (Red → Green → Refactor)
- **Document-First**: 文档优先，任何代码操作前先读取文档
- **智能分流**: 自动判断任务类型 (DOCUMENT / CODE_SIMPLE / CODE_COMPLEX)
- **并行执行**: 关键阶段并行执行

---

## 代码搜索策略

**分层回退**: auggie (语义搜索) → llmdoc (结构化) → Grep (模式匹配)

| 方法 | 精度 | 速度 | 状态 |
|:---|:---:|:---:|:---:|
| **auggie** | 高 | 快 | 推荐 |
| **llmdoc** | 中 | 快 | 默认 |
| **Grep** | 低 | 慢 | 备选 |

### auggie 安装（推荐）
```bash
npm install -g @augmentcode/auggie@prerelease
```

---

## TDD 规范

### 核心流程
```
RED (写测试) → GREEN (写实现) → REFACTOR (重构)
```

### 关键约束
1. 绝不一次性生成代码和测试
2. 先写测试，后写实现
3. 测试必须先失败
4. 覆盖率 > 80%

---


## SpecDD 规范

### 核心流程
```
Phase 4 (Architecture) → Phase 5 (Spec) → Phase 6 (TDD Red)
```

### 适用场景
- 新功能开发
- 跨模块变更 (3+ 文件)
- API 设计
- 数据结构设计

---

## 文件操作规范

| 操作 | 推荐工具 | 禁止使用 |
|:---|:---|:---|
| 读取 | cat, head, tail, `Read` | - |
| 搜索文件 | find, ls, `Glob` | - |
| 搜索内容 | grep, `Grep` | - |
| 创建 | `Write` | touch, echo, cat > |
| 编辑 | `Read` + `Write` | sed, awk, vim |

**原因**: shell 编辑工具容易出错（上下文重复、特殊字符转义问题）

---

## 代码风格

- **不可变性**: 优先使用不可变对象
- **小文件**: 200-400 行，避免 >800 行
- **纯净代码**: 禁止使用 Emoji
- **清晰命名**: 变量/函数名要自解释

---

## Git 操作规范

- **Commit**: Conventional Commits (feat:, fix:, refactor:, docs:, test:)

---

## Agent 调用

Agent 通过**自然语言**调用，无需硬编码。

```
"使用 investigator 调查数据库连接问题"
"让 architect 设计缓存系统架构"
"请 tdd-guide 指导测试编写"
```

---

## 成功指标

1. 测试通过，覆盖率 > 80%
2. 无已知安全漏洞
3. 代码可读性高
4. 需求完整满足
5. 文档已更新（如需要）
