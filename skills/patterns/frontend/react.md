# Frontend Patterns

前端开发通用模式和最佳实践。

## 技能描述

这个技能包含前端开发的通用模式，适用于 React、Vue 等框架。

## 包含的模式

### 1. 组件设计原则

**单一职责**：
- 每个组件只做一件事
- 可组合的小组件优于大组件

**组件分类**：
- UI 组件：纯展示，无业务逻辑
- 容器组件：处理数据和逻辑
- 页面组件：组合多个组件

### 2. 状态管理

**本地状态**：
```typescript
const [count, setCount] = useState(0);
```

**提升状态**：
```typescript
// 父组件管理状态
function Parent() {
  const [value, setValue] = useState('');
  return <Child value={value} onChange={setValue} />;
}
```

**全局状态**（Redux/Pinia）：
```typescript
// 使用 selector 选择需要的切片
const userName = useSelector(state => state.user.name);
```

### 3. 数据获取

**自定义 Hook 模式**：
```typescript
function useFetch<T>(url: string) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    fetch(url)
      .then(res => res.json())
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [url]);

  return { data, loading, error };
}
```

### 4. 表单处理

**受控组件**：
```typescript
const [value, setValue] = useState('');

<input
  value={value}
  onChange={e => setValue(e.target.value)}
/>
```

**表单库（推荐）**：
- React: react-hook-form
- Vue: vee-validate

### 5. 错误边界

```typescript
class ErrorBoundary extends React.Component {
  state = { hasError: false };

  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback />;
    }
    return this.props.children;
  }
}
```

### 6. 性能优化

**React 优化**：
```typescript
// memo 避免不必要渲染
const UserCard = React.memo(({ user }) => <div>{user.name}</div>);

// useMemo 缓存计算结果
const sortedList = useMemo(() => list.sort(), [list]);

// useCallback 缓存函数
const handleClick = useCallback(() => {}, [deps]);
```

**懒加载**：
```typescript
const HeavyComponent = React.lazy(() => import('./Heavy'));

<Suspense fallback={<Spinner />}>
  <HeavyComponent />
</Suspense>
```

### 7. 样式方案

**CSS Modules**：
```css
/* UserCard.module.css */
.container { padding: 16px; }
```

```typescript
import styles from './UserCard.module.css';
<div className={styles.container} />
```

**Tailwind CSS**：
```typescript
<div className="p-4 bg-white rounded-lg shadow" />
```

### 8. API 客户端封装

```typescript
class ApiClient {
  private baseURL: string;

  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`);
    return response.json();
  }

  async post<T>(path: string, data: unknown): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return response.json();
  }
}
```

## 使用场景

- 创建新的前端组件
- 状态管理决策
- 性能优化
- 代码审查参考
