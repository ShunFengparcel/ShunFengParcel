# 查询 EXPLAIN 审核指引

目的：对慢查询路径进行定期审核与调优，保障订单详情与列表接口的稳定性。

## 建议审核项

- 按订单号精确检索：`SELECT * FROM sf_orders WHERE order_no = ? LIMIT 1;`
- 订单审计记录：`SELECT * FROM sf_order_audit_logs WHERE order_id = ? ORDER BY created_at ASC;`
- 改派链路：`SELECT * FROM sf_order_reassignments WHERE order_id = ? ORDER BY created_at ASC;`
- 任务关联：`SELECT * FROM sf_courier_tasks WHERE order_id = ? LIMIT 1;`

## EXPLAIN 示例

```sql
EXPLAIN SELECT * FROM sf_orders WHERE order_no = 'SF202501010001' LIMIT 1;
```

检查要点：

- `type` 至少达到 `ref`/`const`，避免 `ALL` 全表扫描
- `key` 命中 `idx_sf_orders_order_no`
- `rows` 接近 1（单行命中）
- `Extra` 不出现 `Using filesort`/`Using temporary`

## 索引与统计维护

- 定期检查 `information_schema.statistics` 指标，确保索引处于可用状态
- 对高更新表监控 `ANALYZE TABLE` 周期，保持统计信息新鲜
- 大表分区或冷热分离时按分区维度建立局部索引

## 运行时审核（服务侧）

- 服务在订单号精确检索路径附带 `EXPLAIN` 审核（仅日志），无需影响主流程
- 慢查询阈值建议：`> 100ms` 记录 audit log 并告警