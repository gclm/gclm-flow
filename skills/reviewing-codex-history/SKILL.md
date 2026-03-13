---
name: reviewing-codex-history
description: Use when reviewing ~/.codex/history.jsonl to find repeated requests, workflow friction, or candidates for docs, hooks, scripts, or skill updates.
---

# 复盘 Codex 历史

这个 skill 用来消费 `analyze-history.py` 导出的结构化事实，再由模型完成语义梳理、候选经验簇命名和结果路由。它不直接等同于 `updating-domain-skills`。

## 核心规则

- 先运行脚本拿结构化事实，不直接让模型吞原始 `history.jsonl`
- 区分两类结果：工作流/配置改进 vs 领域经验回写
- 只有“领域特有、可复用、已验证”的结论才交给 `updating-domain-skills`
- 工作流、hooks、scripts、README、全局 skills 的改进不要硬塞进领域 skill

## 工作顺序

1. 先运行总览：`python3 ~/.codex/bin/analyze-history.py --output /tmp/codex-history.json`。
2. 读取这个 JSON：先看 `metadata`、`repeated_prompts`、`workflow_topics`、`domain_topics`、`recent_focus`。
3. 如果某个方向证据还不够密，再做定向抽样：
   - 例：`python3 ~/.codex/bin/analyze-history.py --topic devops --topic-kind domain --topic-samples 12 --format markdown`
   - 例：`python3 ~/.codex/bin/analyze-history.py --topic database --topic-kind domain --topic-samples 12 --format json`
4. 基于 `counts + samples + focused_topic` 做语义整理：
   - 给高频领域命名候选经验簇
   - 合并同义问题
   - 识别哪些只是噪音或一次性问题
5. 判断结果归属：
   - workflow/config/tooling：更新 README、hooks、scripts、全局 skills
   - domain-specific：再走 `updating-domain-skills`
6. 只把高频、稳定、值得复用的结论沉淀下来。

## 输出理解

- `analyze-history.py` 只输出事实层：统计、样本、会话摘要、定向 topic 样本。
- 经验簇命名、结论提炼、路由建议属于解释层，应由当前 skill 和模型完成。
- 只有当候选簇里的结论经过真实任务验证后，才进入 `updating-domain-skills`。

## 推荐输出

- `Overview`：高频主题、重复请求、近期焦点
- `Workflow Improvement Candidates`：适合落到 hooks / scripts / README / 全局 skills 的项
- `Domain Writeback Candidates`：适合交给 `updating-domain-skills` 的候选经验
- `Noise / No Action`：不值得沉淀的内容

标准结构见 [output-template.md](references/output-template.md)

## 参考资料

- [action-routing.md](references/action-routing.md)
- [output-template.md](references/output-template.md)

## 联动技能

- `updating-domain-skills`
- `writing-skills`
- `documentation`
