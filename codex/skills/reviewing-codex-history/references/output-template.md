# Output Template

用于统一 `reviewing-codex-history` 的最终输出结构，避免只报统计、不提炼结论。

## 何时查看

- 已经拿到 `analyze-history.py` 导出的 JSON，准备开始做复盘结论输出
- 希望不同会话里的复盘结果保持稳定格式，便于比较和回写

## 模板

```markdown
## Overview
- 时间范围：...
- 高频 workflow 主题：...
- 高频 domain 主题：...
- 近期焦点：...
- 如有定向抽样：补充 `focused_topic` 的关键信号

## Workflow Improvement Candidates
- 主题：...
  证据：引用 `workflow_topics.counts`、`repeated_prompts`、`workflow_topics.samples`
  结论：...
  建议去向：README / hooks / bin / 全局 skills

## Domain Writeback Candidates
- 领域：...
  候选经验簇：...
  证据：引用 `domain_topics.counts`、`domain_topics.samples`，必要时补 `focused_topic.samples`
  结论：...
  前置条件：必须经过真实任务验证
  建议去向：`updating-domain-skills` -> 对应领域 skill / references

## Noise / No Action
- 主题：...
  原因：一次性问题 / 证据不足 / 只是控制词 / 与长期治理无关

## Next Actions
1. 立即可做的 workflow 改进
2. 需要继续验证的 domain 候选
3. 暂不处理的项
```

## 重点做法

- 先写 `Overview`，再分 workflow / domain / noise 三类
- 每个候选项都要带证据来源，至少引用一个 count 和一个 sample
- 定向抽样过的主题，优先引用 `focused_topic.samples`，不要只引用总览统计
- `Domain Writeback Candidates` 不要直接写成最终 skill 内容，先写候选经验与验证前提
- `Workflow Improvement Candidates` 优先输出能减少重复提问或重复操作的项

## 注意事项

- 不要把 `Topic Counts` 直接当结论
- 不要把 workflow 改进误路由到 `updating-domain-skills`
- 没有足够证据时，宁可放进 `Noise / No Action` 或 `待验证`，不要过度提炼
