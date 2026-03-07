---
name: reviewing-codex-history
description: Use when reviewing ~/.codex/history.jsonl to find repeated requests, workflow friction, or candidates for docs, hooks, scripts, or skill updates.
---

# 复盘 Codex 历史

这个 skill 用来把 `~/.codex/history.jsonl` 里的真实交互整理成改进候选。它不直接等同于 `updating-domain-skills`。

## 核心规则

- 先做事实归纳：重复请求、长会话、主题集中度、近期焦点
- 区分两类结果：工作流/配置改进 vs 领域经验回写
- 只有“领域特有、可复用、已验证”的结论才交给 `updating-domain-skills`
- 工作流、hooks、scripts、README、全局 skills 的改进不要硬塞进领域 skill

## 工作顺序

1. 运行 `python3 ~/.codex/bin/analyze-history.py` 或源码目录下的同名脚本。
2. 读报告：先看重复请求、主题分布、近期焦点。
3. 判断结果归属：
   - workflow/config/tooling：更新 README、hooks、scripts、全局 skills
   - domain-specific：再走 `updating-domain-skills`
4. 只把高频、稳定、值得复用的结论沉淀下来。

## 参考资料

- [action-routing.md](references/action-routing.md)

## 联动技能

- `updating-domain-skills`
- `writing-skills`
- `documentation`
