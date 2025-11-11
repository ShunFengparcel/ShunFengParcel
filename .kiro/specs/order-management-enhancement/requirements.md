# 订单管理深化系统需求文档

## 简介

本文档定义顺丰速运订单管理系统的5个深化功能模块，旨在提升订单查询性能、支持大批量数据处理、建立自动化补偿机制、构建风控体系、优化费用计算引擎。

## 术语表

- **System**: 顺丰速运订单管理系统
- **User**: 使用系统的客服人员、运营人员或管理员
- **Order**: 快递订单，包含寄件人、收件人、物流状态等信息
- **ES**: Elasticsearch搜索引擎
- **Compensation**: 订单异常补偿，包括超时、破损、丢失等场景
- **Risk Score**: 风险评分，用于识别异常订单的量化指标
- **Fee Calculator**: 费用计算引擎，根据规则计算订单费用

## 需求

### 需求 1：订单搜索引擎优化

**用户故事：** 作为客服人员，我希望能够快速搜索海量历史订单，以便及时响应客户查询

#### 验收标准

1. WHEN User输入订单号进行精确查询, THE System SHALL在50毫秒内返回查询结果
2. WHEN User使用多个条件组合查询订单, THE System SHALL在500毫秒内返回查询结果
3. WHEN User输入手机号进行模糊查询, THE System SHALL在200毫秒内返回匹配的订单列表
4. WHEN User输入地址关键词进行全文检索, THE System SHALL在300毫秒内返回相关订单并高亮匹配内容
5. THE System SHALL支持至少10亿条订单数据的查询
6. WHEN Order数据在MySQL中发生变更, THE System SHALL在3秒内同步到ES索引
7. THE System SHALL提供订单号、手机号、姓名、地址、状态、时间范围、金额范围等多维度筛选条件

### 需求 2：订单导出与批量处理

**用户故事：** 作为运营人员，我希望能够导出大批量订单数据并进行批量操作，以便进行数据分析和业务处理

#### 验收标准

1. THE System SHALL支持单次导出最多100万条订单数据
2. WHEN User发起导出任务, THE System SHALL以异步方式处理并返回任务ID
3. THE System SHALL以每秒1000条的速度处理导出数据
4. WHEN 导出任务完成, THE System SHALL生成Excel文件并上传到对象存储
5. WHEN 导出任务完成, THE System SHALL通过站内消息通知User
6. THE System SHALL为导出文件生成有效期为7天的下载链接
7. THE System SHALL记录导出任务的状态、进度、文件URL等信息
8. THE System SHALL支持批量打印最多100个订单
9. THE System SHALL支持批量标记最多1000个订单
10. THE System SHALL支持批量分配最多500个订单给快递员

### 需求 3：订单补偿机制

**用户故事：** 作为客服人员，我希望系统能够自动计算补偿金额并触发审批流程，以便快速处理客户投诉

#### 验收标准

1. WHEN Order发生超时配送异常, THE System SHALL自动计算补偿金额为运费金额
2. WHEN Order发生包裹破损异常且已保价, THE System SHALL自动计算补偿金额为保价金额
3. WHEN Order发生包裹破损异常且未保价, THE System SHALL自动计算补偿金额为商品价值与300元的较小值
4. WHEN Order发生包裹丢失异常且已保价, THE System SHALL自动计算补偿金额为保价金额
5. WHEN Order发生包裹丢失异常且未保价, THE System SHALL自动计算补偿金额为商品价值与300元的较小值
6. WHEN 补偿金额小于等于1000元, THE System SHALL自动审批通过
7. WHEN 补偿金额大于1000元, THE System SHALL提交人工审批流程
8. WHEN 补偿审批通过, THE System SHALL在24小时内完成打款
9. WHEN 补偿打款完成, THE System SHALL通过短信和站内消息通知User
10. THE System SHALL记录补偿申请、审批、打款的完整流程日志

### 需求 4：订单风控系统

**用户故事：** 作为风控人员，我希望系统能够实时识别异常订单并自动拦截，以便防止欺诈和刷单行为

#### 验收标准

1. WHEN Order的收寄件地址相同, THE System SHALL标记为高风险并拦截订单
2. WHEN Order的地址无法通过地图API验证, THE System SHALL标记为高风险并拦截订单
3. WHEN User在1小时内创建超过10个订单, THE System SHALL标记为中风险并提交人工审核
4. WHEN User的订单取消率超过50%, THE System SHALL标记为中风险并提交人工审核
5. WHEN 同一寄件人和收件人在1个月内交易超过20次, THE System SHALL标记为高风险并拦截订单
6. WHEN Order的保价金额超过商品价值的2倍, THE System SHALL标记为高风险并提交人工审核
7. THE System SHALL在100毫秒内完成风险评分计算
8. THE System SHALL根据用户信用、订单特征、行为特征计算0到100分的风险评分
9. WHEN 风险评分为0到30分, THE System SHALL标记为低风险并正常流转
10. WHEN 风险评分为31到70分, THE System SHALL标记为中风险并提交人工审核
11. WHEN 风险评分为71到100分, THE System SHALL标记为高风险并拦截订单
12. THE System SHALL记录风控拦截和审核的完整日志

### 需求 5：费用计算引擎优化

**用户故事：** 作为系统开发人员，我希望费用计算引擎能够快速准确地计算订单费用，以便提升用户体验

#### 验收标准

1. THE System SHALL在20毫秒内完成订单费用计算
2. WHEN Order的起始地和目的地在同一省份, THE System SHALL使用首重8元、续重1元每公斤的费率
3. WHEN Order的起始地和目的地跨省且非偏远地区, THE System SHALL使用首重12元、续重2元每公斤的费率
4. WHEN Order的起始地或目的地为偏远地区, THE System SHALL使用首重18元、续重3元每公斤的费率
5. WHEN Order选择保价服务, THE System SHALL计算保价费为保价金额的0.5%且最低2元
6. WHEN Order选择代收货款服务, THE System SHALL计算代收费为代收金额的1%且最低5元
7. WHEN Order在高峰期创建, THE System SHALL在基础费用上加价20%
8. WHEN Order在夜间22点到次日6点创建, THE System SHALL在基础费用上加价30%
9. THE System SHALL将计算结果缓存1小时
10. THE System SHALL实现缓存命中率达到80%以上
11. THE System SHALL支持100个以上的计费规则
12. WHEN 计费规则发生变更, THE System SHALL在不重启服务的情况下热更新规则
