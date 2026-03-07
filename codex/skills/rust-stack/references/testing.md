# Rust Testing Notes

通用测试策略见 `testing`。这里仅保留 Rust 特有细节。

- 优先 `cargo test`
- 需要 HTTP 集成测试时，优先用框架测试工具或最小可运行服务实例
- async 测试明确 runtime、资源初始化和清理
- 尽量断言实际结果与完整对象，而不是零散字段
