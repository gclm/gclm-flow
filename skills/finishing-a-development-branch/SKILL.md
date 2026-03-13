---
name: finishing-a-development-branch
description: Use when implementation appears complete, verification is done, and you need to decide whether to merge, open a PR, keep the branch, or discard the work.
---

# 完成开发分支

把“做完代码”收束为一个可执行的收尾闭环：验证、选路径、执行、清理。

## 核心规则

- 没有最新验证证据，不进入收尾决策。
- 合并、推送、删除分支、删除 worktree 之前，先确认 base branch、当前分支、工作区状态。
- 丢弃工作属于高风险动作，必须显式二次确认。

## 收尾顺序

1. 执行 `verification-before-completion`：确认测试、构建、diff 与目标一致。
2. 确定基线：`main`/`master`/目标集成分支，以及当前 feature branch。
3. 向用户呈现 4 个选项：
   - 本地合并回基线分支
   - 推送并创建 PR
   - 保留当前分支和 worktree
   - 丢弃当前工作
4. 按所选路径执行，并再次验证结果。
5. 如果使用了 worktree，按选项决定是否清理。

## 四个选项

### 1. 本地合并

- 切回基线分支
- 更新基线分支
- 合并当前分支
- 对合并结果再跑一次关键验证
- 成功后可删除 feature branch 与对应 worktree

### 2. 推送并创建 PR

- 推送当前分支到远端
- 基于实际改动生成 PR 标题与摘要
- 写清验证方式
- 保留分支；worktree 可保留直到 PR 生命周期结束

### 3. 保留当前分支

- 不合并、不删除
- 明确告知分支名、worktree 路径、当前状态

### 4. 丢弃当前工作

执行前必须明确：
- 将删除哪个分支
- 将删除哪个 worktree
- 最坏情况是什么
- 如何从 commit 或备份恢复

然后等待用户明确确认。

## Worktree 清理规则

- 本地合并完成后：通常可以清理 worktree
- PR 路径：默认保留，除非用户明确要求清理
- 保留分支：保留 worktree
- 丢弃工作：确认后删除 worktree 与分支

## 联动技能

- `verification-before-completion`
- `using-git-worktrees`
- `code-review`
