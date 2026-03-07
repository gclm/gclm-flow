---
name: go-stack
description: Use when working on Go services, Gin or Echo applications, concurrency, package boundaries, HTTP middleware, or Go-specific build and runtime concerns.
---

# Go 技术栈

这个 skill 负责 Go 项目的入口判断和工程经验索引。主文档只保留适用范围、关键关注点和 references 入口。

## 核心规则

- 先判断问题属于 HTTP 层、并发、数据访问、项目结构还是构建运行
- 通用测试策略看 `testing`；Go 特有测试细节看 `references/testing.md`
- 真实经验优先沉淀到 `references/`

## 重点关注

- package 边界与依赖方向
- handler / service / repository 分层
- context 传递、并发安全、错误处理
- Gin / Echo 中间件与绑定校验

## 参考资料

- [gin.md](references/gin.md)
- [testing.md](references/testing.md)
