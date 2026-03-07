# Action Routing

用于判断 history 复盘结果应该落到哪里，而不是把所有结论都写进领域 skill。

## 何时查看

- 已经跑完 `history.jsonl` 复盘，但还不确定结论该写回哪里
- 同时出现 workflow 改进和领域经验，不想混写

## 重点做法

- 如果结论是 hooks、scripts、README、默认配置、全局 skill 的问题，归到 workflow/config 治理
- 如果结论是 Go、Python、Rust、前端、DevOps、数据库等领域里的稳定经验，归到 `updating-domain-skills`
- 如果只是一次性聊天噪音、偶发疑问、未验证猜测，不回写

## 路由清单

- `README` / `AGENTS`：维护约定、使用说明、配置边界
- `hooks` / `bin`：高频提醒、轻量自动化、发布与校验链路
- 全局 skills：`testing`、`code-review`、`documentation`、`reviewing-codex-history`
- 领域 skills：`devops`、`frontend-stack`、`python-stack`、`go-stack`、`java-stack`、`rust-stack`、`database`

## 注意事项

- `updating-domain-skills` 不是总入口，它只负责领域经验回写
- 复盘结果应优先减少重复提问和重复人工操作，而不是生成更多文档噪音
