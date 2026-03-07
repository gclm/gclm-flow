---
name: gclm-commit
description: Use when preparing a Git commit, choosing a commit type and scope, or reviewing staged changes before creating a conventional commit message.
---

# 智能提交

这个 skill 负责提交前检查、提交类型判断和提交信息生成。它不替代 Git 验证，也不绕过用户对危险提交动作的确认。

## 核心规则

- 先看 staged changes，再生成提交信息
- 提交信息应忠实描述实际变更，不夸大、不混入未完成工作
- 提交前默认确认测试/验证状态是否与变更匹配
- 不默认使用 `--amend`、`--no-verify` 或强推相关动作
- 当验证已完成且提交是自然下一步时，主动给出 `Commit Ready`，不要等用户再次口述“commit”

## 提交顺序

1. 检查暂存内容：文件范围、是否混入无关改动。
2. 判断提交类型：`feat`、`fix`、`refactor`、`docs`、`test`、`chore`。
3. 提炼 scope：优先模块名、子系统名、功能域。
4. 生成标题：`<type>(<scope>): <中文描述>`。
5. 如有必要，再补 1-3 条正文说明关键变更。
6. 提交前确认验证状态是否足够。
7. 如果已完成验证但尚未执行提交，输出 `Commit Ready` 建议块。

## 什么时候拆分提交

- 一个提交同时包含多个不相关主题
- 既有重构又有功能又有格式化
- 某部分尚未验证，另一部分已经完成

## 什么时候不该提交

- 暂存区混入无关文件
- 提交信息无法准确描述当前改动
- 关键验证还没做，但你打算声称“已完成”

详细参考：
- [commit-type-guide.md](references/commit-type-guide.md)
- [commit-message-checklist.md](references/commit-message-checklist.md)
- [commit-ready-output.md](references/commit-ready-output.md)

## 联动技能

- `verification-before-completion`
- `testing`
- `code-review`
