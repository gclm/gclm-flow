---
description: "启动 gclm-flow 融合工作流 - TDD-First + llmdoc 优先 + 多 Agent 并行"
---

# /gclm 命令

启动融合开发工作流。

## 使用方法

```
/gclm <功能描述>
```

## 示例

```
/gclm 实现用户登录功能
/gclm 添加订单导出到 CSV
/gclm 优化数据库查询性能
```

## 9 阶段流程

| 阶段 | 名称 | 说明 |
|:---|:---|:---|
| 0 | llmdoc 读取 | 强制优先读取文档 |
| 1 | Discovery | 需求理解 |
| 2 | Exploration | 并行探索代码库 (3 个 investigator) |
| 3 | Clarification | 强制澄清疑问 |
| 4 | Architecture | 架构设计方案 (2 个 architect) |
| 5 | Spec | 规范文档 (SpecDD) |
| 6 | TDD Red | 编写测试（先失败） |
| 7 | TDD Green | 编写实现 |
| 8 | Refactor+Doc | 重构和更新文档 |
| 9 | Summary | 完成总结 |

## 状态管理

- 状态文件: `.claude/gclm.{task_id}.local.md`
- 中途可恢复
- Stop Hook 保护

## 并行执行

以下阶段并行执行：
- **Phase 2**: 3 个 investigator
- **Phase 4**: 2 个 architect + 1 个 investigator
- **Phase 8**: worker + code-reviewer

## 关键约束

1. **Phase 0**: 强制读取 llmdoc
2. **Phase 3**: 不可跳过，必须澄清
3. **Phase 6**: TDD Red，测试必须先失败
4. **Phase 8**: 必须询问是否更新文档

## 完成信号

```
<promise>GCLM_WORKFLOW_COMPLETE</promise>
```

## 手动退出

如需强制退出工作流：

```bash
sed -i.bak 's/^active: true/active: false/' .claude/gclm.*.local.md
```
