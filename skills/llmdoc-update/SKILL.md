---
name: llmdoc-update
description: 更新项目 llmdoc 文档，基于代码变更同步更新 LLM 优化的项目文档
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash"]
---

# llmdoc-update Skill

更新项目的 llmdoc 文档，保持代码与文档同步。

## 触发时机

- Phase 7 完成后，用户确认需要更新文档
- 重大功能开发完成后
- 架构变更后

## 更新流程

### 1. 扫描变更
```bash
# 获取最近的代码变更
git diff HEAD~1 --name-only
```

### 2. 识别影响范围
- 新增的模块/组件
- 修改的 API
- 变更的架构
- 更新的依赖

### 3. 更新对应文档

#### llmdoc/index.md
- 更新模块列表
- 添加新功能导航

#### llmdoc/overview/
- 更新 project.md（如果范围变更）
- 更新 structure.md（如果结构变更）

#### llmdoc/architecture/
- 更新相关模块文档
- 添加新模块文档

## 文档模板

### 模块文档模板
```markdown
# {模块名称}

## 概述
{模块的简要描述和目的}

## 职责
- 职责 1
- 职责 2

## 公开接口
### {函数/类名}
```typescript
function signature
```
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

## 更新原则

1. **保持同步**: 代码变更后立即更新
2. **LLM 优化**: 使用 LLM 友好的格式
3. **简洁清晰**: 避免冗余，突出重点
4. **交叉引用**: 使用链接连接相关文档

## 质量检查

- [ ] 所有新模块有文档
- [ ] API 变更已反映
- [ ] 依赖关系正确
- [ ] 示例代码可运行
- [ ] 无过时信息
