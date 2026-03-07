# Query Review Checklist

检查这些高频问题：

- 是否有 `SELECT *`
- 是否缺少分页或 limit
- 是否在循环里发起查询，形成 N+1
- where / join / order by 字段是否有合适索引
- 是否做了无意义的排序、去重、聚合
- 是否把本该在数据库做的过滤搬到应用层
- 查询结果集是否过大，导致内存或网络浪费
