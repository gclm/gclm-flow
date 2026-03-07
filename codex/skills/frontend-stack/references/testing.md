# Frontend Testing Notes

通用测试策略见 `testing`。这里仅保留前端特有细节。

- 单元/组件测试：优先 `Vitest` + `React Testing Library`
- 断言用户可见行为，不要过度依赖内部实现细节
- 异步 UI 测试优先等待真实渲染结果，不用随意 `setTimeout`
- E2E 只保留关键用户路径，不把所有场景都堆到浏览器测试里
