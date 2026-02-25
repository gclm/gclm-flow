# 测试规则

所有项目通用的测试规范。

## 测试原则

### FIRST 原则
- **F**ast：测试要快
- **I**ndependent：测试要独立
- **R**epeatable：测试可重复
- **S**elf-validating：自动验证结果
- **T**imely：及时编写测试

### 测试金字塔
```
      /\
     /  \    E2E 测试（少量）
    /----\
   /      \  集成测试（适量）
  /--------\
 /          \ 单元测试（大量）
/------------\
```

## 单元测试

### 命名规范
```
// Java
void methodName_stateUnderTest_expectedBehavior()

// Python
def test_method_name_state_under_test_expected_behavior():

// Go
func TestMethodName_StateUnderTest_ExpectedBehavior(t *testing.T)
```

### AAA 模式
```python
def test_create_user():
    # Arrange（准备）
    user_data = {"email": "test@example.com", "name": "Test"}

    # Act（执行）
    result = user_service.create(user_data)

    # Assert（断言）
    assert result.email == "test@example.com"
    assert result.name == "Test"
```

### 测试覆盖
- 每个公共方法至少一个测试
- 测试正常路径和边界情况
- 测试错误处理

## 集成测试

### 隔离策略
- 使用测试数据库
- 使用内存数据库（如 H2、SQLite）
- Mock 外部服务
- 清理测试数据

### 示例
```python
@pytest.fixture
def test_db():
    db = TestDatabase()
    yield db
    db.cleanup()

def test_user_repository(test_db):
    repo = UserRepository(test_db)
    user = repo.create({"email": "test@example.com"})
    assert user.id is not None
```

## E2E 测试

### 选择场景
- 关键业务流程
- 用户常用路径
- 高风险操作

### 最佳实践
- 使用稳定的测试数据
- 等待机制而非固定延迟
- 独立的测试环境
- 清理测试痕迹

## 测试数据

### 工厂模式
```python
# 使用 factory_boy
class UserFactory(factory.Factory):
    class Meta:
        model = User

    email = factory.Sequence(lambda n: f'user{n}@example.com')
    name = factory.Faker('name')

# 使用
user = UserFactory.create()
```

### Fixtures
```python
@pytest.fixture
def sample_user():
    return UserFactory.create(email="test@example.com")
```

## 覆盖率

### 目标
| 类型 | 目标覆盖率 |
|------|-----------|
| 核心业务逻辑 | 90% |
| 整体项目 | 80% |
| 工具类 | 70% |

### 排除项
- 生成的代码
- 第三方库
- 简单的 getter/setter

## Mock 和 Stub

### 何时使用
- 外部 API 调用
- 数据库操作（单元测试）
- 时间相关操作
- 随机数生成

### 示例
```python
from unittest.mock import Mock, patch

def test_send_email():
    email_service = Mock()
    email_service.send.return_value = True

    result = email_service.send("test@example.com", "Subject", "Body")

    assert result is True
    email_service.send.assert_called_once()
```

## 持续集成

- 每次提交运行测试
- PR 合并前必须通过测试
- 失败时阻止部署
- 生成覆盖率报告
