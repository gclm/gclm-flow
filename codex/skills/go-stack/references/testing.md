# Go Testing Notes

用于补充 Go 测试的栈特有细节；通用测试策略见 `testing`。

## 何时查看

- 需要测试 Go HTTP 层、并发逻辑或标准库风格代码

## 重点做法

- 优先用标准库 `testing`
- HTTP 层优先用 `httptest`
- 并发相关测试要明确超时、goroutine 生命周期和资源清理
- table-driven tests 适合规则型逻辑，但不要为了形式牺牲可读性

## 注意事项

- 避免为了复用表格把断言逻辑写得比被测代码还难读
