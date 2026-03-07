---
name: using-git-worktrees
description: Use when work needs an isolated branch and workspace, especially for larger changes, risky experiments, or when the current tree should remain untouched.
---

# 使用 Git Worktree

需要隔离时就隔离，不在主工作区硬做高风险改动。

## 使用时机

- 大功能、长任务、实验性改造
- 当前工作区已有未提交改动
- 需要并行处理多个分支任务

## 检查项

1. 先确认 worktree 目录策略。
2. 新建分支并创建独立目录。
3. 跑一次基线检查：构建/测试/最小启动。
4. 明确 worktree 路径、分支名、remote 追踪关系。

## 注意

- project-local worktree 要确保被忽略
- 基线不干净时，不要把后续问题误判成你新引入的
