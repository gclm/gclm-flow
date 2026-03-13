---
name: devops
description: Use when working on deployment, CI/CD, containers, Kubernetes, Terraform, release automation, or environment/runtime delivery concerns.
---

# DevOps

这个 skill 负责部署、交付、基础设施与运行时变更的入口判断。主文档只保留触发场景和检查方向，详细经验放到 `references/`。

## 核心规则

- 先识别是构建、发布、容器、编排、基础设施还是运行时配置问题
- 优先验证交付路径和回滚路径，而不是只看配置是否"写得像对"
- 涉及镜像、架构、runner、权限、网络边界时，默认提高风险等级
- 领域经验优先沉淀到 `references/`，不要把长示例继续堆回主文档

## 常见场景

- Docker / 镜像构建
- Kubernetes / 编排与部署
- CI/CD / GitHub Actions / 发布流水线
- Terraform / 基础设施变更
- 运行时配置、Secrets、环境差异
- 代理分组、服务路由、部署选址

## 参考资料

- [github-actions-vm-image-ci.md](references/github-actions-vm-image-ci.md)
- [github-actions-log-debugging.md](references/github-actions-log-debugging.md)
- [deployment-checklist.md](references/deployment-checklist.md)
- [service-routing-and-deployment-patterns.md](references/service-routing-and-deployment-patterns.md)

## 测试与验证边界

- 通用测试策略看 `testing`
- DevOps 相关只保留交付验证、发布前检查、部署后核验等专项内容
