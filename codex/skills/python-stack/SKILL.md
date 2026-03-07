---
name: python-stack
description: Use when working on Python backend architecture, FastAPI or Flask services, dependency wiring, async behavior, packaging, or Python-specific runtime and framework concerns.
---

# Python 技术栈

这个 skill 负责 Python 服务端项目的入口判断和框架经验索引。主文档保留适用范围与 references 入口，详细实践放在 `references/`。

## 核心规则

- 先判断属于框架层、依赖注入、配置、异步行为、数据模型还是打包运行问题
- 通用测试策略看 `testing`；Python 特有测试细节看 `references/testing.md`
- 真实经验优先写回 `references/`，避免主文档继续堆模板代码

## 重点关注

- FastAPI / Flask 路由与依赖注入
- Pydantic / 配置管理
- 异步边界和阻塞调用
- 项目结构、导入边界、错误处理

## 参考资料

- [fastapi.md](references/fastapi.md)
- [testing.md](references/testing.md)
