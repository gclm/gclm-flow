# llmdoc 文档规则

## 核心原则

**任何代码操作前必须优先读取文档**

---

## llmdoc 结构

```
llmdoc/
├── index.md              # 导航入口 - 永远首先阅读
├── overview/             # "这个项目是什么？" - 必须全部阅读
├── architecture/         # "它是怎么工作的？" - LLM 检索地图
├── guides/               # "如何做 X？" - 分步指南
└── reference/            # "X 的具体细节是什么？" - API 规范、约定
```

---

## 读取优先级 (NON-NEGOTIABLE)

任何代码操作前必须：

1. 检查 `llmdoc/` 是否存在
2. **如果存在**:
   - 读取 `llmdoc/index.md`
   - 读取 `llmdoc/overview/*.md` (全部)
   - 根据任务读取相关 `llmdoc/architecture/*.md`
3. **如果不存在**:
   - **自动生成 llmdoc** (无需用户确认)
   - 使用 `investigator` agent 扫描代码库
   - 生成 `llmdoc/index.md`
   - 生成 `llmdoc/overview/` 基础文档
   - 然后继续读取流程

---

## 目录说明

### index.md
- 导航和概览
- 项目入口
- 永远首先阅读

### overview/
- 项目上下文
- 回答"这是什么项目？"
- **必须全部阅读**

### architecture/
- 系统设计
- LLM 检索地图
- 模块关系图
- 回答"它是怎么工作的？"

### guides/
- 操作指南
- 分步指南
- 回答"如何做 X？"

### reference/
- 详细规范
- API 规范
- 数据模型
- 约定
- 回答"X 的具体细节是什么？"

---

## 文档更新

代码变更后必须询问：

```
AskUserQuestion: "是否使用 recorder agent 更新项目文档？"
```

**仅在用户确认后才调用 recorder agent 更新文档**

---

## 文档质量标准

### 清晰度
- 使用简洁的语言
- 避免歧义
- 提供示例

### 完整性
- 覆盖所有关键模块
- 包含必要的上下文
- 链接相关文档

### 可维护性
- 及时更新
- 保持一致性
- 版本控制

---

## LLM 优化

### 关键文件引用
- 保留关键文件路径
- 标注负责的模块
- 说明依赖关系

### 结构化信息
- 使用表格
- 列表清晰
- 代码注释

### 主题串联
- 通过主题组织
- 跨文档链接
- 索引完善

---

## 示例项目

参考：[TokenRoll/minicc/llmdoc](https://github.com/TokenRollAI/minicc/tree/main/llmdoc)

---

## llmdoc 自动生成规则

当检测到项目没有 `llmdoc/` 目录时，自动执行以下流程：

### 1. 代码库扫描 (investigator agent)

**扫描目标**:
- 项目结构（目录、文件组织）
- 主要模块和组件
- 技术栈（语言、框架、工具）
- 入口文件和关键路径
- 测试文件位置

**输出**:
- 项目基本信息
- 模块清单
- 依赖关系
- 关键文件列表

### 2. 生成基础文档

**必须生成**:
```yaml
llmdoc/
├── index.md              # 项目导航和概览
└── overview/
    ├── project.md        # 项目介绍、目标、范围
    ├── tech-stack.md     # 技术栈清单
    └── structure.md      # 目录结构说明
```

### 3. 文档模板

#### index.md 模板
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

#### overview/project.md 模板
```markdown
# 项目介绍

## 项目目标
{提取自 README 或代码分析}

## 项目范围
{主要功能模块}

## 非目标
{明确不包括的内容}
```

### 4. 自动化约束

- **无需用户确认**: 直接执行生成
- **最小化生成**: 只生成基础文档
- **增量完善**: 后续可在 Phase 7 补充 architecture/ 和 guides/
- **保持简洁**: 避免过度生成无用文档

### 5. 生成后处理

生成完成后：
1. 继续执行 Phase 0 读取流程
2. 输出生成摘要
3. 提示用户可以在 Phase 7 选择完善文档
