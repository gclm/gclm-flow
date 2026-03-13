# GitHub Actions Log Debugging with `gh`

用于通过 `gh` CLI 拉取和分析 GitHub Actions 日志，替代浏览器访问 Actions 页面。

## 何时查看

- 用户提供了 Actions run/job URL 需要定位失败原因
- 需要在终端直接读取日志而不打开浏览器


当用户提供一个 GitHub Actions 的链接（run URL、job URL 或 step URL）时，优先用 `gh` CLI 拉取日志，不要直接打开浏览器或调用 chrome-devtools / exa 等网络工具。

## 优先级规则

1. **`gh` CLI 优先**：只要能从 URL 解析出 run/job ID，就用 `gh` 读取日志。
2. **chrome-devtools 仅用于 UI 验证**：Actions 日志是纯文本数据，不需要浏览器渲染，chrome-devtools 在此场景无增量价值。
3. **exa/web search 是最后手段**：当 `gh` 无法访问（权限不足、私有 repo 未授权）时，说明原因后再决定是否走其他路径。

## 从 URL 提取 run ID

GitHub Actions URL 格式：
```
https://github.com/<owner>/<repo>/actions/runs/<run_id>
https://github.com/<owner>/<repo>/actions/runs/<run_id>/jobs/<job_id>
```

提取示例：
```bash
RUN_ID=12345678
REPO="owner/repo"
```

## 常用 `gh` 命令

```bash
# 查看 run 概览（所有 job 状态）
gh run view $RUN_ID --repo $REPO

# 查看完整日志（所有 job）
gh run view $RUN_ID --repo $REPO --log

# 只看失败的步骤日志
gh run view $RUN_ID --repo $REPO --log-failed

# 查看特定 job 的日志
gh run view $RUN_ID --repo $REPO --job $JOB_ID --log

# 列出某个 run 下的所有 job
gh run view $RUN_ID --repo $REPO --json jobs

# 重新触发失败的 run
gh run rerun $RUN_ID --repo $REPO --failed
```

## 诊断顺序

1. 用 `gh run view $RUN_ID --repo $REPO` 先看整体状态和失败 job。
2. 用 `--log-failed` 快速定位错误文本，避免输出全量日志。
3. 根据错误文本定位到具体 step，再用 `--job $JOB_ID --log` 精读上下文。
4. 如需重跑，用 `gh run rerun --failed` 只跑失败部分。

## 适用边界

- 适用于有 `gh` CLI 且已完成 `gh auth login` 的环境。
- 对私有 repo，确认 token 具备 `repo` 和 `workflow` 权限。
- 日志超长时，用 `grep` 或 `| head -n 100` 过滤；不要把全量日志贴进上下文。
- 不适用于需要查看 Actions UI 交互操作（如手动触发 input 表单）的场景，那才需要 chrome-devtools。

## 注意事项

- 日志超长时用 `--log-failed` 或 `grep` 过滤，避免把全量日志贴进上下文。
- 私有 repo 确认 token 具备 `repo` 和 `workflow` 权限。
- 不适用于需要查看 Actions UI 交互操作的场景。
