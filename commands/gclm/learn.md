# /gclm:learn - 记忆管理

管理长期记忆，包括错误记录、模式提取和知识查询。

## 用法

```
/gclm:learn <子命令> [参数]
```

## 子命令

| 子命令 | 描述 |
|--------|------|
| `error` | 记录错误和解决方案 |
| `pattern` | 从代码中提取模式 |
| `search <关键词>` | 搜索记忆 |
| `list` | 列出所有记忆 |
| `stats` | 记忆统计 |
| `clean` | 清理重复/过时记忆 |
| `export` | 导出记忆 |
| `import <文件>` | 导入记忆 |

## 功能

1. **错误记录**
   - 记录遇到的错误
   - 保存解决方案
   - 追踪发生频率

2. **模式提取**
   - 提取成功模式
   - 记录使用场景
   - 建立索引

3. **知识查询**
   - 搜索相关记忆
   - 按标签/语言过滤
   - 相似度排序

## 存储位置

```
~/.gclm-flow/memory/
├── errors/           # 错误记忆
│   ├── java.json
│   ├── python.json
│   ├── go.json
│   ├── rust.json
│   └── frontend.json
├── patterns/         # 模式记忆
│   ├── api-design.json
│   ├── error-handling.json
│   ├── testing.json
│   └── security.json
└── index.json        # 索引文件
```

## 工作流程

1. 调用 `remember` 代理执行操作
2. 更新记忆索引
3. 生成操作报告

## 输出

### 记录结果
```markdown
**已记录到记忆库**

- 类型：错误
- ID：err-20260225-001
- 标签：import, circular-dependency, fastapi
- 相关记忆：2 条
```

### 搜索结果
```markdown
**找到 3 条相关记忆**

#### 相关错误 (2 条)
- [err-001] 循环导入问题 (发生 3 次)
  - 解决方案：将导入语句移到函数内部

#### 相关模式 (1 条)
- [pat-001] 延迟导入模式
  - 使用场景：解决循环依赖
```

### 统计信息
```markdown
# 记忆统计

## 总览
- 总记忆数：50
- 错误记忆：30
- 模式记忆：20

## 按语言分布
- Java: 10 错误, 5 模式
- Python: 12 错误, 8 模式
- Go: 8 错误, 7 模式

## 最近 7 天
- 新增：5 条
- 查询：23 次

## 重复错误 TOP 3
1. 循环导入 (5 次)
2. 空指针异常 (3 次)
3. 类型转换错误 (2 次)
```

## 示例

```bash
# 记录错误
/gclm:learn error

# 提取模式
/gclm:learn pattern src/services/

# 搜索记忆
/gclm:learn search "循环导入"

# 按标签搜索
/gclm:learn search --tag security

# 按语言搜索
/gclm:learn search --lang python

# 查看统计
/gclm:learn stats

# 清理过时记忆
/gclm:learn clean --older-than 90d

# 导出记忆
/gclm:learn export > backup.json

# 导入记忆
/gclm:learn import backup.json
```

## 相关命令

- `/gclm:fix` - 修复问题
- `/gclm:doc` - 文档管理
