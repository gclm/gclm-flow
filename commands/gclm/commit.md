# /gclm:commit - 智能提交

生成规范的 Git 提交信息并执行提交。

## 用法

```
/gclm:commit [选项]
```

## 功能

1. **分析变更**
   - 扫描暂存的文件
   - 分析代码差异
   - 识别变更类型

2. **生成提交信息**
   - 遵循 Conventional Commits 规范
   - 自动推断类型（feat/fix/refactor等）
   - **使用中文**描述变更内容
   - 格式：`<type>(<scope>): <中文描述>`

3. **执行提交**
   - 检查提交前置条件
   - 执行 git commit
   - 可选推送到远程

## 提交类型

| 类型 | 描述 | 示例 |
|------|------|------|
| `feat` | 新功能 | feat: 添加用户登录功能 |
| `fix` | Bug 修复 | fix: 修复登录验证逻辑 |
| `refactor` | 代码重构 | refactor: 优化用户服务 |
| `docs` | 文档更新 | docs: 更新 API 文档 |
| `test` | 测试相关 | test: 添加用户服务测试 |
| `chore` | 杂项 | chore: 更新依赖版本 |
| `style` | 代码风格 | style: 格式化代码 |
| `perf` | 性能优化 | perf: 优化查询性能 |

## 提交信息规范

**格式：**
```
<type>(<scope>): <中文描述>

<中文正文说明>

```

**规则：**
- 标题和正文使用**中文**
- 类型（feat/fix 等）保持英文
- 正文每行以 `-` 开头，说明具体变更

## 工作流程

1. 检查 Git 状态
2. 分析暂存的变更
3. 生成提交信息
4. 确认并执行提交

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

## 选项

- `--push`: 提交后推送到远程
- `--no-verify`: 跳过 pre-commit 钩子
- `--amend`: 修改上次提交
- `--dry-run`: 只显示提交信息，不执行

## 示例

```bash
# 智能提交
/gclm:commit

# 提交并推送
/gclm:commit --push

# 预览提交信息
/gclm:commit --dry-run

# 修改上次提交
/gclm:commit --amend
```

## 相关命令

- `/gclm:review` - 代码审查
- `/gclm:test` - 运行测试
