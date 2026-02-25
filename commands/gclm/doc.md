# /gclm:doc - 文档管理

管理和更新项目文档。

## 用法

```
/gclm:doc <子命令> [选项]
```

## 子命令

| 子命令 | 描述 |
|--------|------|
| `update` | 更新文档以反映代码变更 |
| `generate` | 生成 API 文档 |
| `sync` | 同步文档和代码 |
| `check` | 检查文档完整性 |

## 功能

1. **更新 llmdoc**
   - overview.md - 项目概述
   - guides/ - 使用指南
   - architecture/ - 架构文档
   - reference/ - 参考资料

2. **生成 API 文档**
   - 从代码注释生成
   - OpenAPI/Swagger 支持

3. **记录决策**
   - ADR (架构决策记录)
   - 技术选型记录

## 工作流程

1. 分析代码变更
2. 调用 `recorder` 代理更新文档
3. 验证文档一致性

## 输出

```markdown
# 文档更新报告

## 更新的文件
- llmdoc/overview.md
- llmdoc/guides/getting-started.md
- llmdoc/architecture/api-design.md

## 变更摘要
- 新增 API 端点: 3 个
- 更新配置说明
- 添加部署指南

## 待更新
- [ ] API 文档需要更新
- [ ] 架构图需要同步
```

## 示例

```bash
# 更新所有文档
/gclm:doc update

# 生成 API 文档
/gclm:doc generate api

# 检查文档完整性
/gclm:doc check

# 记录技术决策
/gclm:doc decision "使用 Redis 作为缓存"
```

## 相关命令

- `/gclm:init` - 初始化项目
- `/gclm:learn` - 记忆管理
