---
name: tdd-workflow
description: "测试驱动开发 - 编写新功能、修复 bug 或重构代码时使用。强制 Red-Green-Refactor 循环，80%+ 覆盖率"
---

# TDD Workflow Skill

测试驱动开发流程。

## 触发条件

- 编写新功能
- 修复 bug
- 重构代码
- 添加 API 端点
- 创建新组件

## 核心原则

### 1. 测试优先
**总是先写测试，然后实现代码使测试通过。**

### 2. 覆盖率要求
- 最低 80% 覆盖率（单元 + 集成 + E2E）
- 所有边缘情况覆盖
- 错误场景测试
- 边界条件验证

### 3. 测试类型

#### 单元测试
- 单个函数和工具
- 组件逻辑
- 纯函数
- 助手和工具函数

#### 集成测试
- API 端点
- 数据库操作
- 服务交互
- 外部 API 调用

#### E2E 测试 (Playwright)
- 关键用户流程
- 完整工作流
- 浏览器自动化
- UI 交互

## TDD 工作流步骤

### Step 1: 编写用户故事
```
作为 [角色]，我想要 [操作]，以便 [好处]

示例：
作为用户，我想要语义搜索市场，
以便即使没有确切关键词也能找到相关市场。
```

### Step 2: 生成测试用例
为每个用户故事创建全面测试：

```typescript
describe('Feature Name', () => {
  it('returns expected results for valid input', async () => {
    // 测试实现
  })

  it('handles edge case: empty input', async () => {
    // 测试边缘情况
  })

  it('falls back gracefully when service unavailable', async () => {
    // 测试回退行为
  })

  it('sorts results by relevance score', async () => {
    // 测试排序逻辑
  })
})
```

### Step 3: 运行测试（应该失败）
```bash
npm test
# 测试应该失败 - 我们还没实现
```

### Step 4: 实现代码
编写最小代码使测试通过：

```typescript
// 由测试指导的实现
export async function featureName(input: InputType) {
  // 实现
}
```

### Step 5: 再次运行测试
```bash
npm test
# 测试现在应该通过
```

### Step 6: 重构
在保持测试绿色的同时改进代码质量：
- 消除重复
- 改进命名
- 优化性能
- 增强可读性

### Step 7: 验证覆盖率
```bash
npm run test:coverage
# 验证 80%+ 覆盖率达成
```

## 测试模式

### 单元测试模式 (Jest/Vitest)
```typescript
import { render, screen, fireEvent } from '@testing-library/react'
import { Button } from './Button'

describe('Button Component', () => {
  it('renders with correct text', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByText('Click me')).toBeInTheDocument()
  })

  it('calls onClick when clicked', () => {
    const handleClick = jest.fn()
    render(<Button onClick={handleClick}>Click</Button>)

    fireEvent.click(screen.getByRole('button'))

    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })
})
```

### API 集成测试模式
```typescript
import { NextRequest } from 'next/server'
import { GET } from './route'

describe('GET /api/resource', () => {
  it('returns data successfully', async () => {
    const request = new NextRequest('http://localhost/api/resource')
    const response = await GET(request)
    const data = await response.json()

    expect(response.status).toBe(200)
    expect(data.success).toBe(true)
    expect(Array.isArray(data.data)).toBe(true)
  })

  it('validates query parameters', async () => {
    const request = new NextRequest('http://localhost/api/resource?limit=invalid')
    const response = await GET(request)

    expect(response.status).toBe(400)
  })

  it('handles database errors gracefully', async () => {
    // Mock database failure
    const request = new NextRequest('http://localhost/api/resource')
    // 测试错误处理
  })
})
```

### E2E 测试模式 (Playwright)
```typescript
import { test, expect } from '@playwright/test'

test('user can complete workflow', async ({ page }) => {
  // 导航到页面
  await page.goto('/')
  await page.click('a[href="/target"]')

  // 验证页面加载
  await expect(page.locator('h1')).toContainText('Title')

  // 填写表单
  await page.fill('input[name="field"]', 'value')

  // 提交表单
  await page.click('button[type="submit"]')

  // 验证成功
  await expect(page.locator('text=Success')).toBeVisible()
})
```

## 测试文件组织

```
src/
├── components/
│   ├── Button/
│   │   ├── Button.tsx
│   │   └── Button.test.tsx          # 单元测试
├── app/
│   └── api/
│       └── resource/
│           ├── route.ts
│           └── route.test.ts         # 集成测试
└── e2e/
    └── workflow.spec.ts               # E2E 测试
```

## Mock 外部服务

### Supabase Mock
```typescript
jest.mock('@/lib/supabase', () => ({
  supabase: {
    from: jest.fn(() => ({
      select: jest.fn(() => ({
        eq: jest.fn(() => Promise.resolve({
          data: mockData,
          error: null
        }))
      }))
    }))
  }
}))
```

### Redis Mock
```typescript
jest.mock('@/lib/redis', () => ({
  searchByVector: jest.fn(() => Promise.resolve([
    { id: 'test-1', score: 0.95 }
  ]))
}))
```

### OpenAI Mock
```typescript
jest.mock('@/lib/openai', () => ({
  generateEmbedding: jest.fn(() => Promise.resolve(
    new Array(1536).fill(0.1)
  ))
}))
```

## 测试覆盖验证

### 运行覆盖率报告
```bash
npm run test:coverage
```

### 覆盖率阈值
```json
{
  "coverageThresholds": {
    "global": {
      "branches": 80,
      "functions": 80,
      "lines": 80,
      "statements": 80
    }
  }
}
```

## 常见测试错误

### ❌ 错误：测试实现细节
```typescript
// 不要测试内部状态
expect(component.state.count).toBe(5)
```

### ✅ 正确：测试用户可见行为
```typescript
// 测试用户看到的
expect(screen.getByText('Count: 5')).toBeInTheDocument()
```

### ❌ 错误：脆弱选择器
```typescript
// 容易破坏
await page.click('.css-class-xyz')
```

### ✅ 正确：语义选择器
```typescript
// 对变更稳健
await page.click('button:has-text("Submit")')
await page.click('[data-testid="submit-button"]')
```

### ❌ 错误：无测试隔离
```typescript
// 测试相互依赖
test('creates user', () => { /* ... */ })
test('updates same user', () => { /* 依赖前一个测试 */ })
```

### ✅ 正确：独立测试
```typescript
// 每个测试设置自己的数据
test('creates user', () => {
  const user = createTestUser()
  // 测试逻辑
})

test('updates user', () => {
  const user = createTestUser()
  // 更新逻辑
})
```

## 持续测试

### 开发期间监视模式
```bash
npm test -- --watch
# 文件变更时自动运行测试
```

### 提交前钩子
```bash
# 每次提交前运行
npm test && npm run lint
```

### CI/CD 集成
```yaml
# GitHub Actions
- name: Run Tests
  run: npm test -- --coverage
- name: Upload Coverage
  uses: codecov/codecov-action@v3
```

## 最佳实践

1. **先写测试** - 总是 TDD
2. **单断言每测试** - 聚焦单一行为
3. **描述性测试名称** - 解释测试内容
4. **Arrange-Act-Assert** - 清晰的测试结构
5. **Mock 外部依赖** - 隔离单元测试
6. **测试边缘情况** - Null、undefined、空、大值
7. **测试错误路径** - 不只是快乐路径
8. **保持测试快速** - 单元测试 < 50ms
9. **测试后清理** - 无副作用
10. **审查覆盖率报告** - 识别缺口

## 成功指标

- 80%+ 代码覆盖率
- 所有测试通过（绿色）
- 无跳过或禁用的测试
- 快速测试执行（单元测试 < 30s）
- E2E 测试覆盖关键用户流程
- 测试在生产前捕获 bug

---

**记住**: 测试不是可选的。它们是支持自信重构、快速开发和生产可靠性的安全网。
