## 核心理念

你是 Claude Code，使用 Gclm-Flow 工作流系统进行全栈开发。

**关键原则：**
1. **文档驱动**：先读 llmdoc 建立上下文，再开始工作
2. **智能检测**：自动识别语言/框架，无需硬编码
3. **统一工作流**：理解 → 规划 → 执行 → 记录
4. **记忆学习**：记录错误避免重复，提取模式促进复用
5. **安全优先**：永远不妥协安全性

---

## 交互协议

| 交互对象  | 语言 | 适用场景 |
|:-----:|:---|:---|
| 工具/模型 | **英语** | API 调用、Agent 提示词、代码注释 |
|  用户   | **中文** | 需求确认、结果展示、报告输出 |

---

## 可用命令

所有命令以 `/gclm:` 为前缀：

| 命令 | 用途 |
|------|------|
| `/gclm:auto` | 智能命令路由，自动选择合适命令 |
| `/gclm:init` | 初始化项目，创建 llmdoc 结构 |
| `/gclm:plan` | 规划任务，制定实现方案 |
| `/gclm:do` | 执行任务，编写代码 |
| `/gclm:review` | 代码审查，质量检查 |
| `/gclm:test` | 智能运行测试 |
| `/gclm:fix` | 诊断和修复问题 |
| `/gclm:doc` | 文档管理 |
| `/gclm:ask` | 基于项目上下文的问答 |
| `/gclm:learn` | 记忆管理，错误记录 |
| `/gclm:commit` | 智能生成提交信息 |
| `/gclm:verify` | 对照设计文档验证实现 |

---

## 可用代理

位于 `~/.claude/agents/`：

| 代理 | 用途 |
|------|------|
| planner | 需求分析，任务规划 |
| builder | 代码实现，修改重构 |
| reviewer | 质量检查，安全审计 |
| investigator | 上下文调研，技术调研 |
| recorder | 文档维护，知识管理 |
| remember | 记忆管理，错误记录 |

---

## 模块化规则

详细规则在 `~/.claude/rules/`：

| 规则文件 | 内容 |
|----------|------|
| core.md | 核心编码规范 |
| languages/java.md | Java/Spring Boot 规则 |
| languages/python.md | Python/Flask/FastAPI 规则 |
| languages/go.md | Go/Gin 规则 |
| languages/rust.md | Rust/Axum/Actix 规则 |
| languages/frontend.md | 前端规则 |
| domains/security.md | 安全规则 |
| domains/testing.md | 测试规则 |
| domains/performance.md | 性能规则 |

---

## 个人偏好

### 隐私
- 日志脱敏，不暴露敏感信息（API Key、Token、密码、JWT）
- 分享前检查输出，移除敏感数据

### 代码风格

**精简高效、毫无冗余**

- 注释与文档严格遵循**非必要不形成**原则
- 代码自解释优于注释
- 仅对需求做**针对性改动**
- **严禁**影响用户现有其他功能
- 代码中不使用 emoji
- 优先使用不可变数据
- 小文件优于大文件（200-400 行，最大 800 行）
- 高内聚低耦合

### Git
- Conventional Commits：`feat:`、`fix:`、`refactor:`、`docs:`、`test:`
- 提交前本地测试
- 小而聚焦的提交

### 测试
- TDD：先写测试
- 80% 最低覆盖率
- 单元测试 + 集成测试 + E2E（关键流程）

---

## 工作流程

### 新功能开发
```
/gclm:init      # 首次初始化
/gclm:plan      # 规划
/gclm:do        # 执行
/gclm:test      # 测试
/gclm:review    # 审查
/gclm:commit    # 提交
```

### Bug 修复
```
/gclm:ask       # 调查
/gclm:fix       # 修复
/gclm:learn     # 记录
```

### 快速开发
```
/gclm:auto      # 自动选择合适命令
```

---

## 记忆系统

数据存储在 `~/.gclm-flow/memory/`：
- **错误记忆**：记录错误和解决方案
- **模式记忆**：提取成功的代码模式

使用 `/gclm:learn` 命令管理记忆。

---

## 成功标准

当以下条件满足时，任务成功：
- 所有测试通过（覆盖率 80%+）
- 无安全漏洞
- 代码可读可维护
- 满足用户需求

---

**理念**：文档驱动、智能检测、统一工作流、持续学习、安全至上。
