# /gclm:review - 代码审查

执行全面的代码质量检查。

## 用法

```
/gclm:review [目标] [选项]
```

## 功能

1. **代码质量检查**
   - 代码风格
   - 设计模式
   - 可维护性

2. **安全性检查**
   - 输入验证
   - 认证授权
   - 敏感数据处理
   - 常见漏洞

3. **性能检查**
   - 算法效率
   - 资源管理
   - 查询优化

4. **测试检查**
   - 测试覆盖率
   - 边界测试
   - 错误处理

## 智能检测

自动检测项目语言和框架，应用对应的审查规则：
- Java/Spring Boot → Java 最佳实践
- Python/Flask/FastAPI → Python 最佳实践
- Go/Gin → Go 最佳实践
- Rust/Axum/Actix → Rust 最佳实践
- 前端 → 前端最佳实践

## 工作流程

1. 调用检测系统识别技术栈
2. 调用 `reviewer` 代理执行审查
3. 调用 `remember` 代理检查历史错误
4. 生成审查报告

## 选项

- `--security`: 只进行安全审查
- `--performance`: 只进行性能审查
- `--fix`: 自动修复可修复的问题

## 输出

```markdown
# 代码审查报告

## 概述
- 审查文件: X 个
- 发现问题: Y 个

## 问题列表

### 🔴 严重
1. [文件:行号] 问题描述
   - 建议: ...

### 🟡 警告
1. ...

### 🔵 建议
1. ...

## 测试结果
- 覆盖率: Z%
```

## 示例

```bash
# 审查整个项目
/gclm:review

# 审查特定文件
/gclm:review src/services/

# 安全审查
/gclm:review --security

# 自动修复
/gclm:review --fix
```

## 相关命令

- `/gclm:test` - 运行测试
- `/gclm:fix` - 修复问题
