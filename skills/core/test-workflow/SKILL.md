# Test Workflow

测试工作流，TDD 开发流程。

## 何时使用

- TDD 开发
- 运行测试
- 提升测试覆盖率

## 工作流程

### 1. TDD 流程

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Red       │ -> │   Green     │ -> │  Refactor   │
│   写测试    │    │   写代码    │    │   重构      │
└─────────────┘    └─────────────┘    └─────────────┘
```

### 2. 测试类型

| 类型 | 描述 | 位置 |
|------|------|------|
| 单元测试 | 测试单个函数/类 | `tests/unit/` |
| 集成测试 | 测试模块间交互 | `tests/integration/` |
| E2E 测试 | 测试用户流程 | `tests/e2e/` |

### 3. 测试命令映射

| 语言 | 框架 | 命令 |
|------|------|------|
| Java | JUnit | `mvn test` / `./gradlew test` |
| Python | pytest | `pytest` |
| Go | go test | `go test ./...` |
| Rust | cargo test | `cargo test` |
| 前端 | jest/vitest | `npm test` |

## 测试规范

### 命名规范

```
// Java
class UserServiceTest { }
void findById_shouldReturnUser_whenUserExists() { }

// Python
def test_find_by_id_returns_user_when_exists():

// Go
func TestFindById_ReturnsUser(t *testing.T) {
```

### AAA 模式

```python
def test_create_user():
    # Arrange
    user_data = {"email": "test@example.com"}

    # Act
    result = service.create(user_data)

    # Assert
    assert result.email == "test@example.com"
```

## 覆盖率目标

| 类型 | 目标 |
|------|------|
| 整体 | 80% |
| 核心业务 | 90% |
| 工具类 | 70% |

## 相关命令

- `/gclm:test` - 运行测试
- `/gclm:review` - 代码审查
