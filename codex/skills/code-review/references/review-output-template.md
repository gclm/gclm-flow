# Review Output Template

用于统一 findings-first 风格的 review 输出。

## 模板

```markdown
1. [P1] [path/to/file.ext:123] 结论一句话
   影响：
   依据：
   建议：

2. [P2] [path/to/file.ext:45] 结论一句话
   影响：
   依据：
   建议：

## Open Questions
- ...

## Residual Risks
- ...

## Change Summary
- ...
```

## 使用要求

- 先列 findings，再写摘要
- 严重级别要和实际风险匹配
- 没有问题时也要保留 residual risks 和 change summary
