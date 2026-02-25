# Python 测试模式

使用 pytest、TDD 方法论、fixtures、mocking 和覆盖率的 Python 应用测试策略。

## 何时激活

- 编写新的 Python 代码（遵循 TDD）
- 设计 Python 项目测试套件
- 审查 Python 测试覆盖率
- 搭建测试基础设施

## TDD 工作流

### RED-GREEN-REFACTOR 循环

```python
# 步骤 1: 编写失败的测试 (RED)
def test_add_numbers():
    result = add(2, 3)
    assert result == 5

# 步骤 2: 实现最少代码 (GREEN)
def add(a, b):
    return a + b

# 步骤 3: 按需重构
```

## pytest 基础

### 基本测试结构

```python
import pytest

def test_addition():
    """测试基本加法。"""
    assert 2 + 2 == 4

def test_string_uppercase():
    """测试字符串大写。"""
    text = "hello"
    assert text.upper() == "HELLO"
```

### 断言

```python
# 相等性
assert result == expected
assert result != unexpected

# 真值
assert result  # 真值
assert not result  # 假值
assert result is None

# 成员关系
assert item in collection

# 类型检查
assert isinstance(result, str)

# 异常测试
with pytest.raises(ValueError):
    raise ValueError("错误消息")

# 检查异常消息
with pytest.raises(ValueError, match="无效输入"):
    raise ValueError("无效输入")
```

## Fixtures

### 基本 Fixture

```python
@pytest.fixture
def sample_data():
    return {"name": "Alice", "age": 30}

def test_sample_data(sample_data):
    assert sample_data["name"] == "Alice"
```

### 带设置/清理的 Fixture

```python
@pytest.fixture
def database():
    # 设置
    db = Database(":memory:")
    db.create_tables()
    yield db  # 提供给测试
    # 清理
    db.close()
```

### Fixture 作用域

```python
# 函数作用域（默认）
@pytest.fixture
def temp_file():
    with open("temp.txt", "w") as f:
        yield f

# 模块作用域
@pytest.fixture(scope="module")
def module_db():
    db = Database(":memory:")
    yield db
    db.close()

# 会话作用域
@pytest.fixture(scope="session")
def shared_resource():
    resource = ExpensiveResource()
    yield resource
    resource.cleanup()
```

## 参数化

```python
@pytest.mark.parametrize("input,expected", [
    ("hello", "HELLO"),
    ("world", "WORLD"),
    ("PyThOn", "PYTHON"),
])
def test_uppercase(input, expected):
    assert input.upper() == expected

@pytest.mark.parametrize("a,b,expected", [
    (2, 3, 5),
    (0, 0, 0),
    (-1, 1, 0),
])
def test_add(a, b, expected):
    assert add(a, b) == expected
```

## Mocking

```python
from unittest.mock import patch, Mock

@patch("mypackage.external_api_call")
def test_with_mock(api_call_mock):
    api_call_mock.return_value = {"status": "success"}
    result = my_function()
    api_call_mock.assert_called_once()
    assert result["status"] == "success"
```

## 异步测试

```python
import pytest

@pytest.mark.asyncio
async def test_async_function():
    result = await async_add(2, 3)
    assert result == 5
```

## 测试组织

```
tests/
├── conftest.py           # 共享 fixtures
├── unit/                 # 单元测试
├── integration/          # 集成测试
└── e2e/                  # 端到端测试
```

## 覆盖率

```bash
pytest --cov=mypackage --cov-report=term-missing --cov-report=html
```

| 代码类型 | 目标 |
|---------|------|
| 关键路径 | 100% |
| 一般代码 | 80%+ |

## 最佳实践

**应该：**
- 遵循 TDD（red-green-refactor）
- 每个测试只测一件事
- 使用描述性名称
- 模拟外部依赖
- 测试边界情况

**不应该：**
- 测试实现细节
- 使用复杂条件
- 在测试间共享状态
- 忽略测试失败
