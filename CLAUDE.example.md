## 核心哲学

**TDD-First + llmdoc 优先 + 多 Agent 并行**

1. **TDD-First**: 坚持测试驱动开发，实现代码前先写测试
2. **llmdoc 优先**: 任何代码操作前必须先读取文档
3. **需求先讨论**: Discovery → Exploration → Clarification
4. **Parallel Execution**: 尽可能并行执行任务
5. **Security-First**: 安全永远第一

---

## 环境上下文

- **语言**: 简体中文
- **系统**: macOS
- **交互**: 强制阻断式确认

---

## Phase 0: llmdoc 优先读取 (NON-NEGOTIABLE)

**任何代码操作前必须：**

1. 检查 `llmdoc/` 是否存在
2. **如果存在**:
   - 读取 `llmdoc/index.md`
   - 读取 `llmdoc/overview/*.md` (全部)
   - 根据任务读取相关 `llmdoc/architecture/*.md`
3. **如果不存在**:
   - **自动生成 llmdoc** (无需用户确认)
   - 使用 `investigator` agent 扫描代码库
   - 生成基础文档 (`index.md` + `overview/`)
   - 然后继续读取流程

### llmdoc 结构

```
llmdoc/
├── index.md              # 导航入口 - 永远首先阅读
├── overview/             # "这个项目是什么？" - 必须全部阅读
├── architecture/         # "它是怎么工作的？" - LLM 检索地图
├── guides/               # "如何做 X？" - 分步指南
└── reference/            # "X 的具体细节是什么？" - API 规范、约定
```

---

## 8 阶段工作流

| 阶段 | 名称 | Agent | 并行 |
|:---|:---|:---|:---:|
| 0 | llmdoc 优先读取 | 主 Agent | - |
| 1 | Discovery | `investigator` | - |
| 2 | Exploration | `investigator` x3 | 是 |
| 3 | Clarification | 主 Agent + AskUser | - |
| 4 | Architecture | `architect` x2 + `investigator` | 是 |
| 5 | TDD Red | `tdd-guide` | - |
| 6 | TDD Green | `worker` | - |
| 7 | Refactor + Doc | `worker` + `code-reviewer` | 是 |
| 8 | Summary | `investigator` | - |

---

## 关键约束

1. **llmdoc 优先**: Phase 0 强制执行，不存在时自动生成
2. **TDD 强制**: Phase 5 必须先写测试
3. **Phase 3 不可跳过**: 必须澄清所有疑问
4. **并行优先**: 能并行的任务必须并行执行
5. **选项式编程**: 使用 AskUserQuestion 展示选项
6. **文档更新询问**: Phase 7 必须询问
7. **状态自动化**: 状态文件更新自动进行，无需确认

---

## Agent 体系

| Agent | 职责 | 模型 |
|:---|:---|:---|
| `investigator` | 探索、分析、总结 | Haiku 4.5 |
| `architect` | 架构设计、方案权衡 | Opus 4.5 |
| `worker` | 执行明确定义的任务 | Sonnet 4.5 |
| `tdd-guide` | TDD 流程指导 | Sonnet 4.5 |
| `code-reviewer` | 代码审查 | Sonnet 4.5 |

---

## TDD 核心流程

### Red-Green-Refactor

```
Red (写测试) → Green (写实现) → Refactor (重构)
```

### TDD 约束

1. 绝不一次性生成代码和测试
2. 先写测试，后写实现
3. 测试必须先失败
4. 覆盖率 > 80%

---

## 可用 Skills

| Skill | 触发词 | 描述 |
|:---|:---|:---|
| `/gclm` | "实现功能"、"开发新功能" | 启动完整工作流 |
| `/investigate` | "什么是"、"X怎么工作"、"分析" | 快速代码库调查 |
| `/tdd` | "写测试"、"TDD" | 测试驱动开发 |

---

## 文件操作规范

| 操作 | 使用工具 | 禁止使用 |
|:---|:---|:---|
| 读取 | `Read` | cat, head, tail |
| 创建 | `Write` | touch, echo, cat > |
| 编辑 | `Edit` | sed, awk |
| 搜索文件 | `Glob` | find, ls |
| 搜索内容 | `Grep` | grep |

---

## Git 操作规范

- **原则**: 只读模式为主
- **Commit**: Conventional Commits (feat:, fix:, refactor:, docs:, test:)

---

## 代码风格

- **不可变性**: 优先使用不可变对象
- **小文件**: 200-400 行，避免 >800 行
- **纯净代码**: 禁止使用 Emoji

---

## 成功指标

1. 测试通过，覆盖率 > 80%
2. 无已知安全漏洞
3. 代码可读性高
4. 需求完整满足
5. 文档已更新（如选择）
