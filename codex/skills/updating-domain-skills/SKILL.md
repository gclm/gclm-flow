---
name: updating-domain-skills
description: Use after finishing work in devops or a language stack when you need to decide whether the new lessons should be written back into domain skills or references.
---

# 回写领域 Skills

把一次真实任务里的有效经验，沉淀回 `devops`、`frontend-stack`、`python-stack`、`go-stack`、`java-stack`、`rust-stack` 的 `SKILL.md` 或 `references/`。

## 核心规则

- 只沉淀可复用、可搜索、会再次遇到的经验
- 通用流程不要写回领域 skill，留给 `testing`、`code-review`、`writing-skills`
- 领域 skill 主文档保持薄，详细经验优先写入 `references/`
- 不写 secrets、账号、机器路径、一次性流水账
- 如果需要代理协助提炼，优先使用 `agents/remember.toml`

## 什么时候回写

适合回写：
- 同类问题重复出现两次以上
- 某个框架/工具有非显而易见的坑
- 某种排障或验证路径明显节省后续时间
- 某个实践在当前代码库里已经被验证有效

不适合回写：
- 纯临时上下文
- 普通常识
- 与具体仓库强绑定的偶发细节

## 回写顺序

1. 判断属于哪个领域 skill。
2. 提炼一句触发提示，必要时补到 `SKILL.md`。
3. 把完整经验写入 `references/*.md`。
4. 检查链接、命名、适用范围和敏感信息。
5. 如需更系统地提炼，可调用 `remember` 代理做结构化整理。

详细模板见 [entry-template.md](references/entry-template.md)
