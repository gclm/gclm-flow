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
```

## 核心步骤

本命令调用 `tdd-guide` agent 执行指导。详细规则见 `agents/tdd-guide.md`。

### SCAFFOLD: 定义接口
```typescript
export function featureName(input: InputType): OutputType {
  throw new Error('Not implemented')
}
```

### RED: 编写测试
```typescript
describe('featureName', () => {
  it('should do something', () => {
    expect(featureName(input)).toBe(expected)
  })
})
```

### 验证失败 → GREEN: 实现代码 → 验证通过 → IMPROVE: 重构

### 覆盖率检查
```bash
npm run test:coverage  # 验证 > 80%
```

## 测试类型

| 类型 | 覆盖内容 | 优先级 |
|:---|:---|:---:|
| 单元测试 | 函数、工具、组件 | 必须 |
| 集成测试 | API、数据库、服务 | 必须 |
| E2E 测试 | 关键用户流程 | 推荐 |

## 边缘情况

- [ ] Null/Undefined
- [ ] 空数组/字符串
- [ ] 无效类型
- [ ] 边界值
- [ ] 错误处理
- [ ] 竞态条件

## 覆盖率要求

- **全局**: 80%+
- **关键代码** (认证、支付、安全): 100%

## 约束

1. 绝不一次性生成代码和测试
2. 先写测试，后写实现
3. 测试必须先失败
4. 覆盖率 > 80%
