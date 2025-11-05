-- 索引守护脚本（MySQL 8.x 建议）
-- 目标：在核心查询路径上建立并维护必要索引，保障订单详情与列表查询在表膨胀后仍稳定输出。

-- sf_orders：按订单号精确检索与状态/更新时间维度
ALTER TABLE `sf_orders` ADD INDEX `idx_sf_orders_order_no` (`order_no`);
ALTER TABLE `sf_orders` ADD INDEX `idx_sf_orders_courier_updated` (`courier_id`,`updated_at`);
ALTER TABLE `sf_orders` ADD INDEX `idx_sf_orders_user_updated` (`user_id`,`updated_at`);

-- sf_order_audit_logs：订单状态变更记录按订单与时间扫描
ALTER TABLE `sf_order_audit_logs` ADD INDEX `idx_audit_order_created` (`order_id`,`created_at`);
ALTER TABLE `sf_order_audit_logs` ADD INDEX `idx_audit_action_type` (`action_type`,`created_at`);

-- sf_courier_tasks：任务关联订单检索
ALTER TABLE `sf_courier_tasks` ADD INDEX `idx_tasks_order` (`order_id`);
ALTER TABLE `sf_courier_tasks` ADD INDEX `idx_tasks_courier_status` (`courier_id`,`task_status`);

-- sf_order_reassignments：改派链路按订单与时间扫描
ALTER TABLE `sf_order_reassignments` ADD INDEX `idx_reassign_order_created` (`order_id`,`created_at`);

-- 注意：MySQL 不支持 CREATE INDEX IF NOT EXISTS，重复执行会报错。
-- 建议在上线前通过 information_schema.statistics 检查是否存在：
-- SELECT INDEX_NAME FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'sf_orders';