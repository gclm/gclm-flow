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
2. 读取 `llmdoc/index.md`
3. 读取 `llmdoc/overview/*.md` (全部)
4. 根据任务读取相关 `llmdoc/architecture/*.md`

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
