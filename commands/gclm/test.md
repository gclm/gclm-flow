# /gclm:test - 运行测试

智能运行项目测试。

## 用法

```
/gclm:test [选项]
```

## 功能

1. **智能检测测试框架**
   - Java: JUnit, TestNG
   - Python: pytest, unittest
   - Go: go test
   - Rust: cargo test
   - 前端: jest, vitest

2. **运行测试**
   - 单元测试
   - 集成测试
   - E2E 测试

3. **生成报告**
   - 测试结果
   - 覆盖率报告
   - 失败分析

## 工作流程

1. 检测项目语言和测试框架
2. 选择正确的测试命令
3. 运行测试
4. 分析结果

## 测试命令映射

| 语言 | 框架 | 命令 |
|------|------|------|
| Java | Maven | `mvn test` |
| Java | Gradle | `./gradlew test` |
| Python | pytest | `pytest` |
| Python | unittest | `python -m unittest` |
| Go | go test | `go test ./...` |
| Rust | cargo test | `cargo test` |
| 前端 | jest | `npm test` |
| 前端 | vitest | `vitest run` |

## 选项

- `--coverage`: 生成覆盖率报告
- `--watch`: 监听模式
- `--e2e`: 运行 E2E 测试
- `--filter <pattern>`: 过滤测试

## 输出

```markdown
# 测试报告

## 测试框架
- 检测到: pytest

## 测试结果
- ✅ 通过: 45
- ❌ 失败: 2
- ⏭️ 跳过: 3
- 📊 覆盖率: 78%

## 失败详情

### test_user_login
- 错误: AssertionError
- 文件: tests/test_auth.py:45
- 原因: 预期状态码 200，实际 401
```

## 示例

```bash
# 运行所有测试
/gclm:test

# 生成覆盖率报告
/gclm:test --coverage

# 运行 E2E 测试
/gclm:test --e2e

# 过滤测试
/gclm:test --filter "auth"
```

## 相关命令

- `/gclm:review` - 代码审查
- `/gclm:fix` - 修复问题
