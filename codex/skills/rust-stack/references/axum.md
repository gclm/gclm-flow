# Axum Notes

用于 Axum 项目的 handler、状态和错误边界提醒。

## 何时查看

- 设计 Axum 路由、handler、状态共享、middleware
- code review 中怀疑提取器、状态、错误类型设计不稳定

## 重点关注

- handler 输入输出是否清晰
- 状态共享是否必要且线程安全
- 错误类型是否统一、可向上透传
- 提取器、middleware、response 结构是否分层合理
