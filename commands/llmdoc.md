---
name: llmdoc
description: LLM 优化的项目文档 - 自动生成或更新 llmdoc
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash"]
---

# /llmdoc - LLM 优化的项目文档

**触发时机**: 代码变更后、Phase 7 审查完成后、架构变更时、首次使用时

---

## 自动行为

**无需用户确认** - 自动检测并执行：

| 情况 | 行为 |
|:---|:---|
| `llmdoc/` 不存在 | 自动生成基础文档 |
| `llmdoc/` 已存在 | 扫描变更并更新文档 |

---

## 1. 自动生成（llmdoc 不存在时）

使用 `investigator` agent 扫描代码库，生成基础文档：

```bash
llmdoc/
├── index.md              # 项目导航和概览
└── overview/
    ├── project.md        # 项目介绍、目标、范围
    ├── tech-stack.md     # 技术栈清单
    └── structure.md      # 目录结构说明
```

### 扫描目标
- 项目结构（目录、文件组织）
- 主要模块和组件
- 技术栈（语言、框架、工具）
- 入口文件和关键路径
- 测试文件位置

### index.md 模板
```markdown
# {项目名称} 文档索引

## 概览
[项目简要描述]

## 快速导航
- [项目介绍](overview/project.md)
- [技术栈](overview/tech-stack.md)
- [目录结构](overview/structure.md)

## 关键模块
{根据扫描结果生成关键模块列表}
```

---

## 2. 更新文档（llmdoc 已存在时）

### 扫描变更
```bash
git diff HEAD~1 --name-only
```

### 识别影响范围
- 新增的模块/组件
- 修改的 API
- 变更的架构
- 更新的依赖

### 更新对应文档

| 文档 | 更新时机 |
|:---|:---|
| `index.md` | 模块列表变更、新功能添加 |
| `overview/project.md` | 项目范围变更 |
| `overview/structure.md` | 目录结构变更 |
| `overview/tech-stack.md` | 依赖技术栈变更 |
| `architecture/*.md` | 模块架构变更、新增模块 |

---

## llmdoc 完整结构

```
llmdoc/
├── index.md              # 导航入口 - 永远首先阅读
├── overview/             # "这个项目是什么？"
│   ├── project.md        # 项目介绍、目标、范围
│   ├── tech-stack.md     # 技术栈清单
│   └── structure.md      # 目录结构说明
├── architecture/         # "它是怎么工作的？"
│   └── {module}.md       # 模块架构文档
├── guides/               # "如何做 X？"
│   └── {guide}.md        # 操作指南
└── reference/            # "X 的具体细节是什么？"
    └── {api}.md          # API 规范
```

---

## 模块文档模板

```markdown
# {模块名称}

## 概述
{模块的简要描述和目的}

## 职责
- 职责 1
- 职责 2

## 公开接口
### {函数/类名}
\`\`\`typescript
function signature
\`\`\`
**参数**: ...
**返回**: ...
**异常**: ...

## 依赖关系
- 依赖: ...
- 被依赖: ...

## 使用示例
\`\`\`typescript
// 示例代码
\`\`\`

## 文件位置
- `path/to/file.ts`
```

---

## 更新原则

1. **保持同步**: 代码变更后立即更新
2. **LLM 优化**: 使用 LLM 友好的格式
3. **简洁清晰**: 避免冗余，突出重点
4. **交叉引用**: 使用链接连接相关文档

---

## 质量检查清单

- [ ] 所有新模块有文档
- [ ] API 变更已反映
- [ ] 依赖关系正确
- [ ] 示例代码可运行
- [ ] 无过时信息
- [ ] index.md 导航完整
