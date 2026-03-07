# Action Routing

用于判断 history 复盘结果应该落到哪里，而不是把所有结论都写进领域 skill。

## 何时查看

- 已经跑完 `history.jsonl` 复盘，但还不确定结论该写回哪里
- 已经拿到 `analyze-history.py` 导出的 JSON，但还不确定如何解释它
- 同时出现 workflow 改进和领域经验，不想混写

## 重点做法

- 如果结论是 hooks、scripts、README、默认配置、全局 skill 的问题，归到 workflow/config 治理
- 如果结论是 Go、Python、Rust、前端、DevOps、数据库等领域里的稳定经验，归到 `updating-domain-skills`
- 如果只是一次性聊天噪音、偶发疑问、未验证猜测，不回写

## 事实层输入

- `metadata`：时间范围、条目数、会话数
- `repeated_prompts`：重复用户需求
- `workflow_topics.counts` / `domain_topics.counts`：高频方向
- `workflow_topics.samples` / `domain_topics.samples`：代表样本
- `focused_topic`：定向抽样后的深挖样本
- `recent_focus`：近阶段重点

## 路由清单

- `README` / `AGENTS`：维护约定、使用说明、配置边界
- `hooks` / `bin`：高频提醒、轻量自动化、发布与校验链路
- 全局 skills：`testing`、`code-review`、`documentation`、`reviewing-codex-history`
- 领域 skills：`devops`、`frontend-stack`、`python-stack`、`go-stack`、`java-stack`、`rust-stack`、`database`

## 注意事项

- `updating-domain-skills` 不是总入口，它只负责领域经验回写
- 复盘结果应优先减少重复提问和重复人工操作，而不是生成更多文档噪音
- 不要把脚本里的 `counts` 当成最终结论；它们只是让模型继续梳理的证据
- 对证据薄弱的主题，先用 `focused_topic` 做定向抽样，再决定是否回写
