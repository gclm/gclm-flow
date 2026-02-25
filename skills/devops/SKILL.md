# DevOps Skills

DevOps 相关技能，包括 Docker、部署、CI/CD 等。

## 何时使用

- 容器化部署
- CI/CD 配置
- 环境管理
- 监控配置

## 包含内容

### 1. Docker

#### Dockerfile 最佳实践

```dockerfile
# 多阶段构建
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
EXPOSE 3000
CMD ["node", "dist/main.js"]
```

#### docker-compose
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://db:5432/app
    depends_on:
      - db
      - redis

  db:
    image: postgres:15
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine

volumes:
  pgdata:
```

### 2. CI/CD

#### GitHub Actions
```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci
      - run: npm test
      - run: npm run build
```

#### 部署流程
```
代码提交 → CI 测试 → 构建镜像 → 推送仓库 → 部署更新
```

### 3. 环境管理

#### 环境变量
```bash
# .env.example
DATABASE_URL=
REDIS_URL=
JWT_SECRET=
API_KEY=
```

#### 配置分层
```
环境变量 > 配置文件 > 默认值
```

### 4. 监控

| 类型 | 工具 |
|------|------|
| 日志 | ELK, Loki |
| 指标 | Prometheus, Grafana |
| 追踪 | Jaeger, Zipkin |
| 告警 | Alertmanager, PagerDuty |

## 相关命令

- `/gclm:review --scope security` - 安全审查
