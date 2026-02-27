---
name: gclm-commit
description: |
  智能提交技能。当用户要求提交、commit、git commit 时自动触发。
  包含：(1) 分析变更 (2) 生成规范提交信息 (3) 执行提交
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - commit
    - git
---

# 智能提交

生成规范的 Git 提交信息并执行提交。

## 触发条件

- 用户说"提交"、"commit"
- 用户完成代码变更后

## 工作流程

### 1. 检查 Git 状态

- 扫描暂存的文件
- 分析代码差异
- 识别变更类型

### 2. 生成提交信息

遵循 Conventional Commits 规范：
- 自动推断类型（feat/fix/refactor 等）
- 使用中文描述变更内容
- 格式：`<type>(<scope>): <中文描述>`

### 3. 确认并提交

- 显示生成的提交信息
- 用户确认
- 执行 git commit

## 提交类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `feat` | 新功能 | feat(auth): 添加用户登录功能 |
| `fix` | Bug 修复 | fix(auth): 修复登录验证逻辑 |
| `refactor` | 代码重构 | refactor(user): 优化用户服务 |
| `docs` | 文档更新 | docs(api): 更新 API 文档 |
| `test` | 测试相关 | test(user): 添加用户服务测试 |
| `chore` | 杂项 | chore(deps): 更新依赖版本 |
| `style` | 代码风格 | style: 格式化代码 |
| `perf` | 性能优化 | perf(db): 优化查询性能 |

## 提交信息格式

```
<type>(<scope>): <中文描述>

- 具体变更说明 1
- 具体变更说明 2

```

## 选项

- `--push`: 提交后推送到远程
- `--no-verify`: 跳过 pre-commit 钩子
- `--amend`: 修改上次提交（不推荐）
- `--dry-run`: 只显示提交信息，不执行

## 输出

```markdown
# 提交准备

## 变更文件
- src/services/auth.py (修改)
- src/api/routes.py (修改)
- tests/test_auth.py (新增)

## 变更分析
- 类型: feat
- 范围: auth
- 描述: 添加用户登录功能

## 生成的提交信息

feat(auth): 添加用户登录功能

- 实现基于 JWT 的用户认证
- 添加登录/登出 API 端点
- 添加单元测试

## 确认提交？
[Y/n]
```

## 注意事项

- 不要使用 `--amend` 修改已推送的提交
- 提交前确保测试通过
- 保持提交粒度适中
