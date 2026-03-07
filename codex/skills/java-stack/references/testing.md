# Java Testing Notes

通用测试策略见 `testing`。这里仅保留 Java 特有细节。

- 单元测试优先 `JUnit 5`
- 依赖替身优先 `Mockito`
- Web 层测试优先框架提供的 slice test（如 `@WebMvcTest`）
- 数据层测试优先最小化真实依赖，避免把整个应用上下文都拉起来
