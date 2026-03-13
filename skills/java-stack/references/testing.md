# Java Testing Notes

用于补充 Java 测试的栈特有细节；通用测试策略见 `testing`。

## 何时查看

- 需要为 Java 服务选择测试框架和测试层次

## 重点做法

- 单元测试优先 `JUnit 5`
- 依赖替身优先 `Mockito`
- Web 层测试优先框架提供的 slice test（如 `@WebMvcTest`）
- 数据层测试优先最小化真实依赖，避免把整个应用上下文都拉起来

## 注意事项

- 不要默认用 `@SpringBootTest` 覆盖所有测试层级，这会让反馈过慢且定位模糊
