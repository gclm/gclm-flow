# TDD 规则

## 核心原则

**测试驱动开发 - Red-Green-Refactor**

```
Red (写测试) → Green (写实现) → Refactor (重构)
```

---

## 绝对规则

1. **绝不一次性生成代码和测试**
2. **先写测试，后写实现**
3. **测试必须先失败 (Red)**
4. **覆盖率 > 80%**

---

## TDD 循环

### Red - 编写失败的测试

```typescript
describe('feature', () => {
  it('should do something', () => {
    expect(featureName(input)).toBe(expected)
  })
})
```

**验证**:
```bash
npm test
# 必须失败
```

### Green - 编写最小实现

```typescript
export function featureName(input) {
  // 最小实现使测试通过
  return expected
}
```

**验证**:
```bash
npm test
# 必须通过
```

### Refactor - 改进代码

- 消除重复
- 改进命名
- 优化性能
- **保持测试绿色**

---

## 测试类型

### 单元测试 (必须)

- 函数和工具
- 组件逻辑
- 纯函数
- 辅助函数

### 集成测试 (必须)

- API 端点
- 数据库操作
- 服务交互
- 外部 API

### E2E 测试 (推荐)

- 关键用户流程
- 完整工作流
- UI 交互

---

## 边缘情况检查清单

- [ ] Null/Undefined
- [ ] 空数组/字符串
- [ ] 无效类型
- [ ] 边界值
- [ ] 错误处理
- [ ] 竞态条件
- [ ] 大数据性能
- [ ] 特殊字符

---

## 覆盖率要求

| 类型 | 要求 |
|:---|:---|
| 全局代码 | 80%+ |
| 认证逻辑 | 100% |
| 支付处理 | 100% |
| 安全相关 | 100% |
| 核心业务 | 100% |

---

## 测试质量检查

- [ ] 所有公共函数有单元测试
- [ ] 所有 API 端点有集成测试
- [ ] 关键流程有 E2E 测试
- [ ] 边缘情况已覆盖
- [ ] 错误路径已测试
- [ ] 外部依赖已 mock
- [ ] 测试相互独立
- [ ] 测试名称描述清晰
- [ ] 断言具体有意义
- [ ] 覆盖率 80%+

---

## 测试反模式

### ❌ 测试实现细节
```typescript
expect(component.state.count).toBe(5)
```

### ✅ 测试用户行为
```typescript
expect(screen.getByText('Count: 5')).toBeInTheDocument()
```

### ❌ 测试相互依赖
```typescript
test('creates user', () => { /* ... */ })
test('updates same user', () => { /* 依赖前一个 */ })
```

### ✅ 测试独立
```typescript
test('creates user', () => {
  const user = createTestUser()
})
test('updates user', () => {
  const user = createTestUser()
})
```

---

## Mock 外部依赖

### Supabase
```typescript
jest.mock('@/lib/supabase', () => ({
  supabase: {
    from: jest.fn(() => ({
      select: jest.fn(() => ({
        eq: jest.fn(() => Promise.resolve({ data, error: null }))
      }))
    }))
  }
}))
```

### Redis
```typescript
jest.mock('@/lib/redis', () => ({
  searchByVector: jest.fn(() => Promise.resolve(results))
}))
```

### OpenAI
```typescript
jest.mock('@/lib/openai', () => ({
  generateEmbedding: jest.fn(() => Promise.resolve(embedding))
}))
```

---

## 持续测试

### 开发期间
```bash
npm test -- --watch
```

### 提交前
```bash
npm test && npm run lint
```

### CI/CD
```yaml
- name: Test
  run: npm test -- --coverage
```

---

## 记住

**没有测试就没有代码。测试不是可选的。**
