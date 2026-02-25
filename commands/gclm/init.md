# /gclm:init - 初始化项目

初始化 Gclm-Flow 工作流，创建 llmdoc 文档结构。

## 用法

```
/gclm:init [选项]
```

## 功能

1. **检测项目信息**
   - 自动检测语言和框架
   - 识别构建工具
   - 检测测试框架

2. **创建 llmdoc 结构**
   ```
   llmdoc/
   ├── overview.md          # 项目概述
   ├── guides/              # 使用指南
   │   ├── getting-started.md
   │   ├── development.md
   │   └── deployment.md
   ├── architecture/        # 架构文档
   │   ├── system-design.md
   │   ├── data-model.md
   │   └── api-design.md
   └── reference/           # 参考资料
       ├── decisions/       # ADR
       ├── conventions/     # 项目约定
       └── patterns/        # 代码模式
   ```

3. **生成项目配置**
   - 创建 CLAUDE.md（如果不存在）
   - 更新 .gitignore

## 工作流程

1. 调用 `investigator` 代理扫描项目
2. 调用 `recorder` 代理创建文档结构
3. 生成初始化报告

## 选项

- `--force`: 强制覆盖已存在的 llmdoc
- `--minimal`: 只创建最小结构

## 示例

```bash
# 初始化项目
/gclm:init

# 强制覆盖
/gclm:init --force

# 最小结构
/gclm:init --minimal
```

## 相关命令

- `/gclm:doc` - 更新文档
- `/gclm:plan` - 规划任务
