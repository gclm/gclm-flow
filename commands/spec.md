---
name: spec
description: SpecDD (Specification-Driven Development) - 规范驱动开发工作流
---

# /spec - SpecDD 规范驱动开发

**触发时机**: 复杂模块开发、需要详细技术规范时

## SpecDD 核心流程

```
Phase 4 (Architecture)
    ↓
Phase 4.5 (Spec) ← 当前
    ↓
Phase 5 (TDD Red)
```

## Spec 文档结构

```markdown
# {功能名称} 规范文档

## 1. 概述
### 1.1 目标
### 1.2 范围
### 1.3 非目标

## 2. 功能需求
### 2.1 用户故事
### 2.2 验收标准
### 2.3 边界条件

## 3. API 设计
### 3.1 公开接口
### 3.2 数据结构
### 3.3 错误处理

## 4. 技术设计
### 4.1 组件架构
### 4.2 数据流
### 4.3 依赖关系

## 5. 测试策略
### 5.1 单元测试覆盖
### 5.2 集成测试场景
### 5.3 边界测试

## 6. 非功能需求
### 6.1 性能要求
### 6.2 安全要求
### 6.3 可维护性
```

## 质量检查清单

- [ ] 目标清晰明确
- [ ] 验收标准可测试
- [ ] API 设计完整
- [ ] 技术方案可行
- [ ] 测试策略全面

## 输出位置

`.claude/specs/{feature-name}.md`

## 与其他阶段关联

- **Phase 1 Discovery**: 提供需求和验收标准
- **Phase 4 Architecture**: 基于设计方案编写
- **Phase 5 TDD Red**: 作为测试编写依据
- **Phase 6 TDD Green**: 作为实现指南
