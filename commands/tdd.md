---
description: "测试驱动开发 - Red-Green-Refactor，80%+ 覆盖率"
---

# /tdd 命令

执行测试驱动开发流程。

## 使用方法

```
/tdd <功能描述>
```

## 示例

```
/tdd 实现用户密码验证
/tdd 添加订单金额计算函数
/tdd 修复登录超时 bug
```

## TDD 循环

```
RED → GREEN → REFACTOR → REPEAT

RED:      编写失败的测试
GREEN:    编写最小实现使测试通过
REFACTOR: 改进代码，保持测试通过
REPEAT:   下一个功能/场景
```

## 流程步骤

### Step 1: 定义接口 (SCAFFOLD)
```typescript
export function featureName(input: InputType): OutputType {
  throw new Error('Not implemented')
}
```

### Step 2: 编写测试 (RED)
```typescript
describe('featureName', () => {
  it('should do something', () => {
    expect(featureName(input)).toBe(expected)
  })
})
```

### Step 3: 运行测试 - 验证失败
```bash
npm test
# 应该失败
```

### Step 4: 实现代码 (GREEN)
```typescript
export function featureName(input: InputType): OutputType {
  // 最小实现
  return expected
}
```

### Step 5: 运行测试 - 验证通过
```bash
npm test
# 应该通过
```

### Step 6: 重构 (IMPROVE)
- 消除重复
- 改进命名
- 优化性能

### Step 7: 检查覆盖率
```bash
npm run test:coverage
# 验证 > 80%
```

## 测试类型

| 类型 | 覆盖内容 | 优先级 |
|:---|:---|:---:|
| 单元测试 | 函数、工具、组件 | 必须 |
| 集成测试 | API、数据库、服务 | 必须 |
| E2E 测试 | 关键用户流程 | 推荐 |

## 边缘情况检查

- [ ] Null/Undefined
- [ ] 空数组/字符串
- [ ] 无效类型
- [ ] 边界值
- [ ] 错误处理
- [ ] 竞态条件

## 覆盖率要求

- **全局**: 80%+
- **关键代码**: 100%
  - 认证逻辑
  - 支付处理
  - 安全相关
  - 核心业务逻辑

## 约束

1. **绝不一次性生成代码和测试**
2. **先写测试，后写实现**
3. **测试必须先失败 (RED)**
4. **覆盖率 > 80%**
