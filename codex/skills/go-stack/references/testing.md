# Go Testing Notes

通用测试策略见 `testing`。这里仅保留 Go 特有细节。

## 何时查看

- 需要测试 Go HTTP 层、并发逻辑或标准库风格代码

## 重点做法

- 优先用标准库 `testing`
- HTTP 层优先用 `httptest`
- 并发相关测试要明确超时、goroutine 生命周期和资源清理
- table-driven tests 适合规则型逻辑，但不要为了形式牺牲可读性
