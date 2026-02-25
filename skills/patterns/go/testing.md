# Go 测试模式

包括表驱动测试、子测试、基准测试、模糊测试和覆盖率的 Go 测试模式。

## 何时激活

- 编写新的 Go 函数或方法
- 为现有代码添加测试覆盖
- 创建基准测试
- 遵循 Go 项目 TDD 工作流

## TDD 工作流

```go
// 步骤 1: 编写失败的测试 (RED)
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}

// 步骤 2: 实现最少代码 (GREEN)
func Add(a, b int) int {
    return a + b
}

// 步骤 3: 按需重构
```

## 表驱动测试

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"正数", 2, 3, 5},
        {"负数", -1, -2, -3},
        {"零值", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

## 测试辅助函数

```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("打开数据库失败: %v", err)
    }

    t.Cleanup(func() {
        db.Close()
    })

    return db
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("意外错误: %v", err)
    }
}
```

## 接口 Mocking

```go
type UserRepository interface {
    GetUser(id string) (*User, error)
}

type MockUserRepository struct {
    GetUserFunc func(id string) (*User, error)
}

func (m *MockUserRepository) GetUser(id string) (*User, error) {
    return m.GetUserFunc(id)
}

func TestUserService(t *testing.T) {
    mock := &MockUserRepository{
        GetUserFunc: func(id string) (*User, error) {
            return &User{ID: "123", Name: "Alice"}, nil
        },
    }

    service := NewUserService(mock)
    user, err := service.GetUserProfile("123")

    assertNoError(t, err)
    if user.Name != "Alice" {
        t.Errorf("got name %q; want %q", user.Name, "Alice")
    }
}
```

## 基准测试

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData(1000)
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// 运行: go test -bench=. -benchmem
```

## 模糊测试

```go
func FuzzParseJSON(f *testing.F) {
    f.Add(`{"name": "test"}`)

    f.Fuzz(func(t *testing.T, input string) {
        var result map[string]interface{}
        err := json.Unmarshal([]byte(input), &result)

        if err != nil {
            return // 无效 JSON 是预期的
        }

        _, err = json.Marshal(result)
        if err != nil {
            t.Errorf("Marshal 失败: %v", err)
        }
    })
}
```

## HTTP Handler 测试

```go
func TestAPIHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("got status %d; want %d", w.Code, http.StatusOK)
    }
}
```

## 覆盖率

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

| 代码类型 | 目标 |
|---------|------|
| 核心逻辑 | 100% |
| 一般代码 | 80%+ |

## 最佳实践

**应该：**
- 先写测试（TDD）
- 使用表驱动测试
- 在辅助函数中使用 `t.Helper()`
- 使用 `t.Cleanup()` 清理
- 使用有意义的测试名称

**不应该：**
- 直接测试私有函数
- 在测试中使用 `time.Sleep()`
- 忽略不稳定的测试
