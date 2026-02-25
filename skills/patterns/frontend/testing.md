# 前端测试模式

使用 Vitest/Jest、React Testing Library 和 Playwright/Cypress 的前端应用测试策略。

## 何时激活

- 编写新的 React/Vue/Angular 组件
- 为前端代码添加测试覆盖
- 搭建 E2E 测试
- 遵循 TDD 工作流

## TDD 工作流

```typescript
// 步骤 1: 编写失败的测试 (RED)
describe('Button', () => {
  it('应该渲染标签', () => {
    render(<Button>点击我</Button>);
    expect(screen.getByText('点击我')).toBeInTheDocument();
  });
});

// 步骤 2: 实现组件 (GREEN)
export function Button({ children }: { children: React.ReactNode }) {
  return <button>{children}</button>;
}

// 步骤 3: 按需重构
```

## 使用 Vitest 进行单元测试

### 基本测试结构

```typescript
import { describe, it, expect, beforeEach } from 'vitest';

describe('Calculator', () => {
  let calculator: Calculator;

  beforeEach(() => {
    calculator = new Calculator();
  });

  it('应该相加两个数字', () => {
    expect(calculator.add(2, 3)).toBe(5);
  });

  it('应该处理负数', () => {
    expect(calculator.add(-1, 1)).toBe(0);
  });
});
```

### 断言

```typescript
// 相等性
expect(value).toBe(expected);
expect(value).toEqual({ name: 'Alice' });

// 真值
expect(value).toBeTruthy();
expect(value).toBeFalsy();
expect(value).toBeNull();
expect(value).toBeUndefined();

// 数字
expect(value).toBeGreaterThan(10);
expect(value).toBeLessThanOrEqual(20);

// 字符串
expect(value).toContain('substring');
expect(value).toMatch(/regex/);

// 数组
expect(array).toHaveLength(3);
expect(array).toContain(item);

// 对象
expect(object).toHaveProperty('name');
expect(object).toMatchObject({ name: 'Alice' });

// 异常
expect(() => fn()).toThrow(Error);
```

## React Testing Library

### 组件测试

```typescript
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('LoginForm', () => {
  it('应该渲染表单字段', () => {
    render(<LoginForm />);

    expect(screen.getByLabelText(/邮箱/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/密码/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /登录/i })).toBeInTheDocument();
  });

  it('应该使用有效数据提交表单', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(<LoginForm onSubmit={onSubmit} />);

    await user.type(screen.getByLabelText(/邮箱/i), 'test@example.com');
    await user.type(screen.getByLabelText(/密码/i), 'password123');
    await user.click(screen.getByRole('button', { name: /登录/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      });
    });
  });

  it('应该显示无效邮箱错误', async () => {
    render(<LoginForm />);

    await userEvent.type(screen.getByLabelText(/邮箱/i), 'invalid');
    await userEvent.click(screen.getByRole('button', { name: /登录/i }));

    expect(await screen.findByText(/无效邮箱/i)).toBeInTheDocument();
  });
});
```

### 测试 Hooks

```typescript
import { renderHook, act } from '@testing-library/react';
import { useCounter } from './useCounter';

describe('useCounter', () => {
  it('应该增加计数', () => {
    const { result } = renderHook(() => useCounter());

    act(() => {
      result.current.increment();
    });

    expect(result.current.count).toBe(1);
  });
});
```

### 测试异步组件

```typescript
describe('UserProfile', () => {
  it('应该显示加载状态然后显示用户数据', async () => {
    render(<UserProfile userId="1" />);

    expect(screen.getByText(/加载中/i)).toBeInTheDocument();

    expect(await screen.findByText('Alice')).toBeInTheDocument();
    expect(screen.queryByText(/加载中/i)).not.toBeInTheDocument();
  });
});
```

## Mocking

### 模拟函数

```typescript
import { vi } from 'vitest';

const mockFn = vi.fn();
mockFn.mockReturnValue('value');
mockFn.mockImplementation((x) => x * 2);

expect(mockFn).toHaveBeenCalled();
expect(mockFn).toHaveBeenCalledWith('arg');
expect(mockFn).toHaveBeenCalledTimes(2);
```

### 模拟模块

```typescript
vi.mock('./api', () => ({
  fetchUser: vi.fn().mockResolvedValue({ name: 'Alice' }),
}));

vi.mock('axios', () => ({
  default: {
    get: vi.fn().mockResolvedValue({ data: {} }),
  },
}));
```

### 模拟定时器

```typescript
beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

it('应该在延迟后调用回调', () => {
  const callback = vi.fn();
  setTimeout(callback, 1000);

  vi.advanceTimersByTime(1000);

  expect(callback).toHaveBeenCalled();
});
```

## 使用 Playwright 进行 E2E 测试

### 基本测试

```typescript
import { test, expect } from '@playwright/test';

test('应该成功登录', async ({ page }) => {
  await page.goto('/login');

  await page.fill('[name="email"]', 'test@example.com');
  await page.fill('[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('h1')).toContainText('欢迎');
});
```

### Page Object 模式

```typescript
class LoginPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto('/login');
  }

  async login(email: string, password: string) {
    await this.page.fill('[name="email"]', email);
    await this.page.fill('[name="password"]', password);
    await this.page.click('button[type="submit"]');
  }
}

test('登录流程', async ({ page }) => {
  const loginPage = new LoginPage(page);
  await loginPage.goto();
  await loginPage.login('test@example.com', 'password123');

  await expect(page).toHaveURL('/dashboard');
});
```

## 测试命令

```bash
# 运行单元测试
vitest

# 运行带覆盖率
vitest --coverage

# 运行 E2E 测试
playwright test

# 运行特定测试文件
vitest Button.test.tsx

# 监视模式
vitest watch

# 更新快照
vitest -u
```

## 覆盖率

```typescript
// vitest.config.ts
export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      exclude: ['node_modules/', 'tests/'],
    },
  },
});
```

| 代码类型 | 目标 |
|---------|------|
| 组件 | 90%+ |
| Hooks | 100% |
| 工具函数 | 80%+ |
| 关键流程 | E2E 覆盖 |

## 最佳实践

**应该：**
- 测试用户交互，而非实现细节
- 使用可访问性查询（getByRole, getByLabelText）
- 使用 `waitFor` 测试异步行为
- 模拟外部依赖
- 使用 userEvent 而非 fireEvent

**不应该：**
- 测试实现细节
- 使用 container.querySelector
- 模拟一切
- 跳过清理
- 忽略可访问性
