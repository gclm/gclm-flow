# Rust Testing Notes

通用测试策略见 `testing`。这里仅保留 Rust 特有细节。

## 何时查看

- 需要为 Rust 服务设计测试入口、async 测试或集成测试

## 重点做法

- 优先 `cargo test`
- 需要 HTTP 集成测试时，优先用框架测试工具或最小可运行服务实例
- async 测试明确 runtime、资源初始化和清理
- 尽量断言实际结果与完整对象，而不是零散字段
