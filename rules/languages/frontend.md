# 前端规则

TypeScript/React/Vue 项目的编码规范和最佳实践。

## TypeScript 编码规范

### 1. 命名规范
```typescript
// 接口/类型：PascalCase
interface UserProfile {}
type UserRole = 'admin' | 'user';

// 类：PascalCase
class UserService {}

// 函数/变量：camelCase
function getUserById(id: number): User {}
const userName = 'test';

// 常量：UPPER_SNAKE_CASE
const MAX_RETRY_COUNT = 3;
const API_BASE_URL = '/api/v1';

// 枚举：PascalCase for type，PascalCase for values
enum UserStatus {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
}

// 私有属性：# 前缀（ES2022+）
class Example {
  #privateField: string;
}

// 文件名：kebab-case
// user-profile.tsx
// api-client.ts
```

### 2. 类型定义
```typescript
// 优先使用 interface（可扩展）
interface User {
  id: number;
  email: string;
  name: string;
  role: UserRole;
}

// 使用 type 用于联合类型、映射类型
type UserRole = 'admin' | 'user' | 'guest';
type UserKeys = keyof User;

// 泛型约束
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

// 工具类型
type PartialUser = Partial<User>;
type ReadonlyUser = Readonly<User>;
type UserWithoutId = Omit<User, 'id'>;
```

### 3. 项目结构
```
src/
├── api/                 # API 调用
│   ├── client.ts
│   └── users.ts
├── components/          # 通用组件
│   ├── ui/
│   └── layout/
├── hooks/               # 自定义 Hooks
│   ├── useAuth.ts
│   └── useUsers.ts
├── pages/               # 页面组件
│   ├── Home.tsx
│   └── Users.tsx
├── store/               # 状态管理
│   ├── index.ts
│   └── userSlice.ts
├── types/               # 类型定义
│   └── index.ts
├── utils/               # 工具函数
│   └── helpers.ts
└── App.tsx
```

## React 最佳实践

### 1. 组件定义
```typescript
// 优先使用函数组件
interface UserCardProps {
  user: User;
  onEdit?: (user: User) => void;
}

export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  return (
    <div className="user-card">
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      {onEdit && <button onClick={() => onEdit(user)}>Edit</button>}
    </div>
  );
};
```

### 2. 状态管理
```typescript
// 使用 useState
const [count, setCount] = useState<number>(0);

// 使用 useReducer 复杂状态
type State = { count: number };
type Action = { type: 'increment' } | { type: 'decrement' };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'increment':
      return { count: state.count + 1 };
    case 'decrement':
      return { count: state.count - 1 };
    default:
      return state;
  }
}

const [state, dispatch] = useReducer(reducer, { count: 0 });

// 使用 Context
const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) throw new Error('useUser must be used within UserProvider');
  return context;
};
```

### 3. 自定义 Hooks
```typescript
// 数据获取 Hook
function useUsers() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    fetchUsers()
      .then(setUsers)
      .catch(setError)
      .finally(() => setLoading(false));
  }, []);

  return { users, loading, error };
}

// 使用
function UserList() {
  const { users, loading, error } = useUsers();

  if (loading) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <ul>
      {users.map(user => <UserCard key={user.id} user={user} />)}
    </ul>
  );
}
```

### 4. 表单处理
```typescript
import { useForm } from 'react-hook-form';

interface LoginForm {
  email: string;
  password: string;
}

function LoginForm() {
  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>();

  const onSubmit = (data: LoginForm) => {
    // 处理登录
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input
        {...register('email', {
          required: 'Email is required',
          pattern: {
            value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
            message: 'Invalid email address'
          }
        })}
        placeholder="Email"
      />
      {errors.email && <span>{errors.email.message}</span>}

      <input
        type="password"
        {...register('password', { required: 'Password is required', minLength: 8 })}
        placeholder="Password"
      />
      {errors.password && <span>{errors.password.message}</span>}

      <button type="submit">Login</button>
    </form>
  );
}
```

### 5. API 调用
```typescript
// API 客户端
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.baseUrl}${path}`);
    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
    return response.json();
  }

  async post<T>(path: string, data: unknown): Promise<T> {
    const response = await fetch(`${this.baseUrl}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
    return response.json();
  }
}

export const api = new ApiClient('/api/v1');

// API 函数
export const userApi = {
  getAll: () => api.get<ApiResponse<User[]>>('/users'),
  getById: (id: number) => api.get<ApiResponse<User>>(`/users/${id}`),
  create: (data: CreateUserRequest) => api.post<ApiResponse<User>>('/users', data),
};
```

## Vue 最佳实践

### 1. 组件定义（Composition API）
```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

interface Props {
  userId: number;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'update', user: User): void;
}>();

const user = ref<User | null>(null);
const loading = ref(true);

const fullName = computed(() => {
  if (!user.value) return '';
  return `${user.value.firstName} ${user.value.lastName}`;
});

onMounted(async () => {
  user.value = await fetchUser(props.userId);
  loading.value = false;
});
</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="user">
    <h2>{{ fullName }}</h2>
    <p>{{ user.email }}</p>
    <button @click="emit('update', user)">Update</button>
  </div>
</template>
```

### 2. 状态管理（Pinia）
```typescript
import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [] as User[],
    currentUserId: null as number | null,
  }),

  getters: {
    currentUser: (state) =>
      state.users.find(u => u.id === state.currentUserId),
  },

  actions: {
    async fetchUsers() {
      const response = await userApi.getAll();
      this.users = response.data;
    },

    setCurrentUser(id: number) {
      this.currentUserId = id;
    },
  },
});
```

## 测试规范

### 1. 单元测试（Vitest）
```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { UserCard } from './UserCard';

describe('UserCard', () => {
  it('renders user information', () => {
    const user: User = {
      id: 1,
      name: 'Test User',
      email: 'test@example.com',
    };

    render(<UserCard user={user} />);

    expect(screen.getByText('Test User')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', async () => {
    const user: User = { id: 1, name: 'Test', email: 'test@example.com' };
    const onEdit = vi.fn();

    render(<UserCard user={user} onEdit={onEdit} />);

    await userEvent.click(screen.getByText('Edit'));

    expect(onEdit).toHaveBeenCalledWith(user);
  });
});
```

### 2. 组件测试
```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import UserCard from './UserCard.vue';

describe('UserCard', () => {
  it('renders user information', () => {
    const wrapper = mount(UserCard, {
      props: {
        user: { id: 1, name: 'Test User', email: 'test@example.com' },
      },
    });

    expect(wrapper.text()).toContain('Test User');
    expect(wrapper.text()).toContain('test@example.com');
  });
});
```

## 性能优化

### 1. React 优化
```typescript
// 使用 memo 避免不必要渲染
const UserCard = React.memo<UserCardProps>(({ user, onEdit }) => {
  return <div>{user.name}</div>;
});

// 使用 useMemo 和 useCallback
function UserList({ users }: { users: User[] }) {
  const sortedUsers = useMemo(() => {
    return [...users].sort((a, b) => a.name.localeCompare(b.name));
  }, [users]);

  const handleEdit = useCallback((user: User) => {
    // 处理编辑
  }, []);

  return sortedUsers.map(user => (
    <UserCard key={user.id} user={user} onEdit={handleEdit} />
  ));
}

// 懒加载
const UserDashboard = React.lazy(() => import('./UserDashboard'));

function App() {
  return (
    <Suspense fallback={<Spinner />}>
      <UserDashboard />
    </Suspense>
  );
}
```

### 2. 代码分割
```typescript
// 动态导入
const routes = [
  {
    path: '/users',
    component: () => import('./pages/Users.vue'),
  },
];
```
