# Memory Skills

记忆系统技能，管理错误记录和模式提取。

## 何时使用

- 记录遇到的错误和解决方案
- 提取可复用的代码模式
- 查询历史知识
- 避免重复犯错

## 记忆类型

| 类型 | 内容 | 存储位置 |
|------|------|----------|
| 错误记忆 | 错误 + 解决方案 | `~/.gclm-flow/memory/errors/` |
| 模式记忆 | 成功的代码模式 | `~/.gclm-flow/memory/patterns/` |

## 工作流程

### 记录错误

```
1. 遇到错误
2. 分析原因
3. 找到解决方案
4. 调用 /gclm:learn error 记录
```

### 提取模式

```
1. 完成任务
2. 识别成功模式
3. 调用 /gclm:learn pattern 提取
```

### 查询记忆

```
1. 开始新任务
2. 调用 /gclm:learn search 查询相关记忆
3. 应用历史知识
```

## 数据格式

### 错误记忆
```json
{
  "id": "err-20260225-001",
  "language": "python",
  "error": {
    "type": "ImportError",
    "message": "cannot import name 'Depends'"
  },
  "solution": {
    "description": "循环导入问题",
    "actions": ["将导入语句移到函数内部"]
  },
  "occurrences": 2,
  "tags": ["import", "circular-dependency"]
}
```

### 模式记忆
```json
{
  "id": "pat-20260225-001",
  "type": "api-design",
  "name": "统一响应格式",
  "description": "所有 API 返回统一的响应结构",
  "when_to_use": ["新建 REST API"],
  "usage_count": 15
}
```

## 相关命令

- `/gclm:learn error` - 记录错误
- `/gclm:learn pattern` - 提取模式
- `/gclm:learn search` - 搜索记忆
- `/gclm:learn stats` - 记忆统计
