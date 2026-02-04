# 命令参考

## 概览

gclm-flow 提供 5 个主要命令，覆盖不同的开发场景：

| 命令 | 用途 | 工作流 |
|:---|:---|:---|
| `/gclm` | 智能分流工作流 | 自动选择 |
| `/investigate` | 代码库调查 | 简化流程 |
| `/tdd` | 测试驱动开发 | TDD 流程 |
| `/spec` | 规范驱动开发 | SpecDD 流程 |
| `/llmdoc` | 文档生成/更新 | 单步操作 |

---

## /gclm

### 语法

```
/gclm <任务描述>
```

### 功能

智能分流工作流，自动判断任务类型并选择最优开发流程。

### 工作流类型

| 类型 | 检测关键词 | 阶段数 |
|:---|:---|:---:|
| 📝 DOCUMENT | 文档、方案、设计、需求 | 7 |
| 🔧 CODE_SIMPLE | bug、修复、error | 6 |
| 🚀 CODE_COMPLEX | 功能、模块、开发 | 9 |

### 示例

```
/gclm 添加用户认证功能
/gclm 修复登录按钮样式
/gclm 编写 API 设计文档
/gclm 重构数据访问层
```

### 阶段流程

```
Phase 0 (llmdoc Reading)
    ↓
Phase 1 (Discovery) - 自动分类
    ↓
Phase 2 (Exploration) - CODE_COMPLEX only
    ↓
Phase 3 (Clarification) - 确认需求
    ↓
Phase 4 (Architecture) - CODE_COMPLEX only
    ↓
Phase 5 (Spec) - CODE_COMPLEX only
    ↓
Phase 6 (TDD Red / Draft)
    ↓
Phase 7 (TDD Green / Refine)
    ↓
Phase 8 (Refactor + Security + Review)
    ↓
Phase 9 (Summary)
```

---

## /investigate

### 语法

```
/investigate <问题>
```

### 功能

快速代码库调查，使用 investigator agent 分析项目。

### 示例

```
/investigate 项目中如何处理用户认证？
/investigate 错误处理机制在哪里？
/investigate 数据库连接是怎么建立的？
/investigate 找出所有 API 端点
```

### 输出

- 相关文件列表 (含行号)
- 调用流程图
- 代码规范识别
- 简洁报告 (< 150 行)

---

## /tdd

### 语法

```
/tdd <功能>
```

### 功能

测试驱动开发，遵循 Red-Green-Refactor 循环。

### TDD 循环

```
Red (写测试) → Green (写实现) → Refactor (重构)
```

### 绝对规则

1. 绝不一次性生成代码和测试
2. 先写测试，后写实现
3. 测试必须先失败
4. 覆盖率 > 80%

### 示例

```
/tdd 添加密码验证函数
/tdd 实现用户注册 API
/tdd 编写数据访问层测试
```

### 阶段

1. **TDD Red**: 编写失败的测试
2. **TDD Green**: 编写最小实现
3. **Refactor**: 重构优化

---

## /spec

### 语法

```
/spec <功能>
```

### 功能

规范驱动开发，先写详细规范文档，再编写测试和实现。

### 适用场景

- 新功能开发
- 跨模块变更 (3+ 文件)
- API 设计
- 数据结构设计

### 示例

```
/spec 设计用户权限系统
/spec 定义支付流程规范
/spec 设计消息队列架构
```

### 阶段

1. **Architecture**: 架构设计
2. **Spec**: 编写规范文档
3. **TDD Red**: 基于规范编写测试
4. **TDD Green**: 实现代码

### 输出文件

`.claude/specs/{feature-name}.md`

---

## /llmdoc

### 语法

```
/llmdoc
```

### 功能

自动生成或更新项目 llmdoc 文档。

### 行为

1. 检查 `llmdoc/` 是否存在
2. 存在 → 扫描代码库并更新文档
3. 不存在 → 生成基础文档

### 生成内容

```
llmdoc/
├── index.md              # 导航入口
├── overview/
│   ├── project.md        # 项目介绍
│   ├── tech-stack.md     # 技术栈
│   └── structure.md      # 目录结构
├── architecture/
│   ├── workflow.md       # 工作流架构
│   ├── agents.md         # Agent 体系
│   └── code-search.md    # 代码搜索策略
├── guides/
│   ├── installation.md   # 安装指南
│   └── quickstart.md     # 快速开始
└── reference/
    ├── commands.md       # 命令参考
    └── configuration.md  # 配置参考
```

---

## 命令对比

| 命令 | 复杂度 | 适用场景 | Agent 使用 |
|:---|:---:|:---|:---|
| `/gclm` | 自动 | 所有场景 | 全部 |
| `/investigate` | 低 | 代码理解 | investigator |
| `/tdd` | 中 | 功能实现 | tdd-guide + worker |
| `/spec` | 高 | 架构设计 | architect + spec-guide + tdd-guide + worker |
| `/llmdoc` | 低 | 文档更新 | investigator |

---

## 最佳实践

### 选择合适的命令

| 需求 | 推荐命令 |
|:---|:---|
| 不确定用哪个 | `/gclm` |
只想了解代码 | `/investigate` |
简单功能开发 | `/tdd` |
复杂架构设计 | `/spec` |
更新项目文档 | `/llmdoc` |

### 组合使用

```
# 1. 先调查代码
/investigate 当前的认证机制

# 2. 然后设计新方案
/spec 添加 OAuth2 认证

# 3. 实现功能
/tdd 实现 OAuth2 客户端
```

---

## 技巧

### 模糊描述

```
/gclm 优化性能
# 会进入 Phase 3 澄清具体优化目标
```

### 精确描述

```
/gclm 添加用户角色管理功能，包含 admin/user/guest 三种角色
# 可以跳过更多阶段
```

### 引用上下文

```
/gclm 像 user service 一样实现 order service
# 可以利用相似功能加速开发
```
