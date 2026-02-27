# AGENTS.md

Gclm-Flow 全局配置，用于 Codex CLI。

## 核心理念

使用 Gclm-Flow 工作流系统进行全栈开发。

**关键原则：**
1. **文档驱动**：先读 llmdoc 建立上下文
2. **智能检测**：自动识别语言/框架
3. **统一工作流**：理解 → 规划 → 执行 → 记录
4. **记忆学习**：记录错误，提取模式
5. **安全优先**：永远不妥协安全性

---

## 交互协议

| 交互对象 | 语言 |
|:-------:|:-----|
| 工具/模型 | 英语 |
| 用户 | 中文 |

---

## 可用技能

| 技能 | 用途 |
|------|------|
| `gclm` | 智能编排（推荐入口） |
| `gclm-init` | 项目初始化 |
| `gclm-commit` | 智能提交 |
| `code-review` | 代码审查 |
| `testing` | 测试 |
| `documentation` | 文档 |
| `memory` | 记忆系统 |

语言栈（自动检测）：`java-stack` `python-stack` `go-stack` `rust-stack` `frontend-stack`

---

## Skills 维护规范

- `SKILL.md` 保持精简入口，聚焦触发条件和最短执行路径
- 详细流程、检查清单、长案例统一放在 `references/`
- 在 `SKILL.md` 中通过链接引用 `references/*.md`，按需展开
- 避免在 `SKILL.md` 持续堆积长排障记录
- 新增内容超过约 30 行时，优先拆分到 `references/`
- 参考文件遵循单主题原则（例如：`github-actions-vm-image-ci.md`）
- 新增经验时采用“双写入”：
- `SKILL.md` 增加一句触发提示
- `references/` 增加完整可复用内容
- 当 `SKILL.md` 明显膨胀时，立即做分层重构

### CI/架构经验（来自 vm-images 事故）

- 工作流中禁止硬编码架构二进制（例如 `yq_linux_amd64`）
- 基于 `uname -m` 做架构映射后再下载二进制
- 设计镜像定制流程时，必须前置校验 host/guest 架构兼容性

---

## 代码风格

- 注释**非必要不形成**
- 代码自解释优于注释
- 仅做**针对性改动**
- 不使用 emoji
- 小文件优于大文件（200-400 行）

---

## Git 规范

```
<type>(<scope>): <中文描述>
```

类型：`feat` `fix` `refactor` `docs` `test` `chore`

---

## 测试规范

- TDD，80% 覆盖率
- 单元 + 集成 + E2E

---

**理念**：文档驱动、智能检测、统一工作流、持续学习、安全至上。
