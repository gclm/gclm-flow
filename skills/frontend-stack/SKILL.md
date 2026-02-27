---
name: frontend-stack
description: |
  前端技术栈完整开发指南。当检测到前端项目（package.json 且无后端框架）
  或用户明确要求 React/Vue/Angular 开发时自动触发。包含：
  (1) 项目结构规范 (2) React/Vue 最佳实践 (3) 测试模式 (4) 状态管理
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - frontend
    - react
    - vue
    - typescript
---

# 前端技术栈开发指南

## 框架检测

- 存在 `react` 依赖 → React，详见 [react.md](references/react.md)
- 存在 `vue` 依赖 → Vue
- 测试相关 → 详见 [testing.md](references/testing.md)

## 标准项目结构

```
src/
├── api/                 # API 调用
├── components/          # 通用组件
│   ├── ui/
│   └── layout/
├── hooks/               # 自定义 Hooks
├── pages/               # 页面组件
├── store/               # 状态管理
├── types/               # 类型定义
├── utils/               # 工具函数
└── App.tsx
```

## 核心规范

### 组件定义 (React)

```typescript
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

### 自定义 Hooks

```typescript
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
```

### API 调用

```typescript
class ApiClient {
  private baseUrl: string;

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
    return response.json();
  }
}
```

### 状态管理 (Zustand)

```typescript
import { create } from 'zustand';

interface UserStore {
  users: User[];
  addUser: (user: User) => void;
}

export const useUserStore = create<UserStore>((set) => ({
  users: [],
  addUser: (user) => set((state) => ({ users: [...state.users, user] })),
}));
```

## 测试规范

- Vitest + React Testing Library
- 目标覆盖率：80%+

## 相关技能

- `code-review` - 前端代码审查
- `testing` - 前端测试模式
