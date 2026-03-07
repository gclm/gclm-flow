---
name: verification-before-completion
description: Use before claiming something is done, fixed, passing, or ready to merge. Requires fresh evidence from commands, diffs, previews, or runtime checks.
---

# 完成前验证

没有证据，不要宣称完成。

## 适用范围

在说下面这些话之前必须先验证：
- 已完成
- 已修复
- 测试通过
- 可以合并
- 没问题

## 验证顺序

1. 明确结论对应的验证动作是什么。
2. 运行完整验证命令或实际预览。
3. 读取输出、退出码、失败数。
4. 结论必须和证据完全一致。

## 例子

- 说“测试通过”前，先跑对应测试命令
- 说“已同步到 `~/.codex`”前，先跑 diff 或实际检查文件
- 说“页面正常”前，先实际打开/预览，而不是只看日志
