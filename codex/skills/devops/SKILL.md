---
name: devops
description: |
  DevOps 技能。当用户要求部署、docker、kubernetes、ci/cd、terraform 时自动触发。
  包含：(1) Docker (2) Kubernetes (3) CI/CD (4) Terraform
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - devops
    - docker
    - kubernetes
    - ci-cd
---

# DevOps

## 触发条件

- 用户要求部署、Docker、Kubernetes、CI/CD、Terraform
- GitHub Actions 中构建 VM 镜像失败
- `virt-customize/libguestfs` 在 CI 里报错

## 经验入口

- GitHub Actions + VM 镜像构建排障：见 [github-actions-vm-image-ci.md](references/github-actions-vm-image-ci.md)
- 涉及多架构构建时，优先校验 runner 与 guest 架构匹配
- 工作流中避免硬编码架构二进制，统一按 `uname -m` 映射下载

## Docker

### Dockerfile 最佳实践

```dockerfile
# 多阶段构建
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
CMD ["node", "dist/main.js"]
```

### docker-compose

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://db:5432/mydb
    depends_on:
      - db
  db:
    image: postgres:15
    volumes:
      - pgdata:/var/lib/postgresql/data
```

## Kubernetes

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 3000
```

## CI/CD

### GitHub Actions

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: npm ci
      - run: npm test
```

## Terraform

```hcl
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"
  tags = {
    Name = "example"
  }
}
```
