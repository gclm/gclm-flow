# Packaging Guide

## 什么时候拆文件

保留在 `SKILL.md`：
- 触发条件
- 核心步骤
- 输出要求
- 红线与边界

拆到 `references/`：
- 长检查清单
- 大段案例
- 领域专项指南
- 输出模板集合

拆到 `scripts/`：
- 反复执行的命令序列
- 模板初始化脚本
- 自动校验或生成动作

## 推荐结构

```text
skill-name/
  SKILL.md
  references/
    checklist.md
    examples.md
  scripts/
    verify.sh
```

## 设计目标

- 主文档短
- 触发明确
- 扩展内容按需加载
- 未来维护成本低
