# Python Testing Notes

通用测试策略见 `testing`。这里仅保留 Python 特有细节。

- 优先 `pytest`
- Web 层测试优先走框架提供的 test client，而不是一上来就启完整服务
- 涉及 async 时，明确 event loop 和 async fixture 边界
- mock 只隔离外部依赖，不要把核心业务逻辑全 mock 掉
