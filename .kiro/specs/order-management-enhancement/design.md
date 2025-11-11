# 订单管理深化系统设计文档

## 概述

本文档描述顺丰速运订单管理深化系统的技术设计方案，包括5个核心模块：订单搜索引擎、订单导出批量处理、订单补偿机制、订单风控系统、费用计算引擎优化。

系统基于现有的Go + Kratos微服务架构，采用MySQL作为主数据库，Redis作为缓存，新增Elasticsearch作为搜索引擎，消息队列用于异步任务处理。

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      前端层                              │
│  微信小程序 / Web管理后台 / 快递员APP                   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    API网关层                             │
│  Kratos HTTP Server (端口: 18000)                       │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    服务层                                │
│  ┌──────────┬──────────┬──────────┬──────────┐         │
│  │订单搜索  │导出服务  │补偿服务  │风控服务  │         │
│  │服务      │          │          │          │         │
│  └──────────┴──────────┴──────────┴──────────┘         │
│  ┌──────────────────────────────────────────┐           │
│  │        费用计算服务                       │           │
│  └──────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    数据层                                │
│  ┌──────────┬──────────┬──────────┬──────────┐         │
│  │MySQL     │Redis     │ES        │OSS       │         │
│  │主数据库  │缓存      │搜索引擎  │文件存储  │         │
│  └──────────┴──────────┴──────────┴──────────┘         │
└─────────────────────────────────────────────────────────┘
```

### 技术栈选型

| 组件 | 技术选型 | 版本 | 用途 |
|------|----------|------|------|
| 后端框架 | Kratos | v2.8.0 | 微服务框架 |
| 编程语言 | Go | 1.21+ | 服务端开发 |
| 数据库 | MySQL | 8.0+ | 主数据存储 |
| 缓存 | Redis | 6.0+ | 缓存和分布式锁 |
| 搜索引擎 | Elasticsearch | 8.0+ | 订单全文检索 |
| 消息队列 | Redis Stream | 6.0+ | 异步任务队列 |
| 对象存储 | 阿里云OSS | - | 导出文件存储 |
| ORM | GORM | v1.31.0 | 数据库操作 |
| 数据同步 | Canal | 1.1.7 | MySQL到ES同步 |

## 模块设计

### 1. 订单搜索引擎模块

#### 1.1 组件设计

```
internal/
├── biz/
│   └── search.go              # 搜索业务逻辑
├── data/
│   ├── elasticsearch.go       # ES客户端封装
│   └── search_repo.go         # 搜索数据仓库
└── service/
    └── search_service.go      # 搜索服务接口
```

#### 1.2 Elasticsearch索引设计

**索引名称**: `sf_orders`

**索引映射**:
```json
{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "ik_max_word": {
          "type": "custom",
          "tokenizer": "ik_max_word"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "order_id": {"type": "long"},
      "order_no": {"type": "keyword"},
      "user_id": {"type": "long"},
      "sender_name": {
        "type": "text",
        "analyzer": "ik_max_word",
        "fields": {"keyword": {"type": "keyword"}}
      },
      "sender_phone": {"type": "keyword"},
      "sender_address": {"type": "text", "analyzer": "ik_max_word"},
      "receiver_name": {
        "type": "text",
        "analyzer": "ik_max_word",
        "fields": {"keyword": {"type": "keyword"}}
      },
      "receiver_phone": {"type": "keyword"},
      "receiver_address": {"type": "text", "analyzer": "ik_max_word"},
      "order_status": {"type": "keyword"},
      "order_amount": {"type": "double"},
      "product_type": {"type": "keyword"},
      "payment_method": {"type": "keyword"},
      "created_at": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
      "updated_at": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss"}
    }
  }
}
```

#### 1.3 数据同步方案

**方案选择**: Canal + Redis Stream

**同步流程**:
1. Canal监听MySQL Binlog
2. 解析INSERT/UPDATE/DELETE事件
3. 发送到Redis Stream
4. 消费者从Stream读取并同步到ES
5. 记录同步状态和失败重试

**同步延迟**: < 3秒

#### 1.4 搜索API设计

**接口**: `POST /api/v1/orders/search`

**请求参数**:
```go
type SearchOrderRequest struct {
    Keyword     string    `json:"keyword"`      // 关键词（订单号/手机号/姓名）
    StartDate   string    `json:"start_date"`   // 开始日期
    EndDate     string    `json:"end_date"`     // 结束日期
    Status      []string  `json:"status"`       // 订单状态列表
    ProductType string    `json:"product_type"` // 产品类型
    MinAmount   float64   `json:"min_amount"`   // 最小金额
    MaxAmount   float64   `json:"max_amount"`   // 最大金额
    Page        int32     `json:"page"`         // 页码
    PageSize    int32     `json:"page_size"`    // 每页数量
}
```

**响应结果**:
```go
type SearchOrderResponse struct {
    Total   int64          `json:"total"`    // 总数
    List    []*OrderInfo   `json:"list"`     // 订单列表
    Took    int64          `json:"took"`     // 查询耗时(ms)
}
```

### 2. 订单导出与批量处理模块

#### 2.1 组件设计

```
internal/
├── biz/
│   ├── export.go              # 导出业务逻辑
│   └── batch.go               # 批量操作逻辑
├── data/
│   ├── export_repo.go         # 导出任务仓库
│   └── oss.go                 # OSS客户端封装
└── service/
    ├── export_service.go      # 导出服务
    └── batch_service.go       # 批量操作服务
```

#### 2.2 导出任务表设计

```sql
CREATE TABLE sf_export_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id VARCHAR(32) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    task_name VARCHAR(200),
    query_condition TEXT,
    total_count INT DEFAULT 0,
    processed_count INT DEFAULT 0,
    file_url VARCHAR(500),
    status VARCHAR(20) NOT NULL,
    progress INT DEFAULT 0,
    error_msg TEXT,
    created_at DATETIME NOT NULL,
    completed_at DATETIME,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 2.3 异步导出架构

**任务队列**: Redis Stream

**Worker数量**: 5个并发Worker

**处理流程**:
1. 用户发起导出请求
2. 创建导出任务记录
3. 任务ID加入Redis Stream
4. Worker消费任务
5. 分批查询数据（每批1000条）
6. 生成Excel文件
7. 上传到OSS
8. 更新任务状态
9. 通知用户

#### 2.4 导出API设计

**创建导出任务**: `POST /api/v1/orders/export`

**请求参数**:
```go
type CreateExportRequest struct {
    TaskName   string              `json:"task_name"`
    Condition  SearchOrderRequest  `json:"condition"`
    Fields     []string            `json:"fields"`
}
```

**响应结果**:
```go
type CreateExportResponse struct {
    TaskId string `json:"task_id"`
}
```

**查询导出任务**: `GET /api/v1/orders/export/{task_id}`

**响应结果**:
```go
type ExportTaskInfo struct {
    TaskId         string `json:"task_id"`
    TaskName       string `json:"task_name"`
    Status         string `json:"status"`
    Progress       int32  `json:"progress"`
    TotalCount     int32  `json:"total_count"`
    ProcessedCount int32  `json:"processed_count"`
    FileUrl        string `json:"file_url"`
    CreatedAt      string `json:"created_at"`
    CompletedAt    string `json:"completed_at"`
}
```

### 3. 订单补偿机制模块

#### 3.1 组件设计

```
internal/
├── biz/
│   ├── compensation.go        # 补偿业务逻辑
│   └── compensation_rule.go   # 补偿规则引擎
├── data/
│   └── compensation_repo.go   # 补偿数据仓库
└── service/
    └── compensation_service.go # 补偿服务
```

#### 3.2 补偿表设计

```sql
CREATE TABLE sf_compensations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    compensation_id VARCHAR(32) UNIQUE NOT NULL,
    order_id BIGINT NOT NULL,
    order_no VARCHAR(20) NOT NULL,
    incident_type VARCHAR(20) NOT NULL,
    compensation_amount DECIMAL(10,2) NOT NULL,
    compensation_reason TEXT,
    status VARCHAR(20) NOT NULL,
    need_approval BOOLEAN DEFAULT FALSE,
    approver_id BIGINT,
    approval_time DATETIME,
    payment_status VARCHAR(20),
    payment_time DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_order_id (order_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 3.3 补偿规则引擎

**规则定义**:
```go
type CompensationRule struct {
    IncidentType string
    Calculator   func(*Order, *Incident) *Compensation
}

var CompensationRules = map[string]CompensationRule{
    "TIMEOUT": {
        IncidentType: "TIMEOUT",
        Calculator: func(order *Order, incident *Incident) *Compensation {
            return &Compensation{
                Amount:       order.FreightFee,
                Reason:       "超时配送，免除运费",
                NeedApproval: false,
            }
        },
    },
    "DAMAGE": {
        IncidentType: "DAMAGE",
        Calculator: func(order *Order, incident *Incident) *Compensation {
            amount := 300.0
            if order.InsuredValue > 0 {
                amount = order.InsuredValue
            } else if order.GoodsValue < 300 {
                amount = order.GoodsValue
            }
            return &Compensation{
                Amount:       amount,
                Reason:       "包裹破损，按保价赔付",
                NeedApproval: amount > 1000,
            }
        },
    },
}
```

#### 3.4 补偿API设计

**创建补偿申请**: `POST /api/v1/compensations`

**请求参数**:
```go
type CreateCompensationRequest struct {
    OrderId      int64  `json:"order_id"`
    IncidentType string `json:"incident_type"`
    Description  string `json:"description"`
    Evidence     []string `json:"evidence"`
}
```

**审批补偿**: `POST /api/v1/compensations/{id}/approve`

**请求参数**:
```go
type ApproveCompensationRequest struct {
    Approved       bool    `json:"approved"`
    AdjustedAmount float64 `json:"adjusted_amount"`
    Comment        string  `json:"comment"`
}
```

### 4. 订单风控系统模块

#### 4.1 组件设计

```
internal/
├── biz/
│   ├── risk.go                # 风控业务逻辑
│   └── risk_rule.go           # 风控规则引擎
├── data/
│   └── risk_repo.go           # 风控数据仓库
└── service/
    └── risk_service.go        # 风控服务
```

#### 4.2 风控记录表设计

```sql
CREATE TABLE sf_risk_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    order_no VARCHAR(20) NOT NULL,
    user_id BIGINT NOT NULL,
    risk_score INT NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    risk_reasons TEXT,
    action VARCHAR(20) NOT NULL,
    reviewer_id BIGINT,
    review_result VARCHAR(20),
    review_time DATETIME,
    created_at DATETIME NOT NULL,
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_risk_level (risk_level),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 4.3 风控规则引擎

**评分模型**:
```go
type RiskScorer struct {
    Rules []RiskRule
}

type RiskRule struct {
    Name      string
    Weight    float64
    Calculate func(*Order, *UserProfile) float64
}

func (rs *RiskScorer) CalculateScore(order *Order, profile *UserProfile) int {
    totalScore := 0.0
    for _, rule := range rs.Rules {
        score := rule.Calculate(order, profile)
        totalScore += score * rule.Weight
    }
    return int(totalScore)
}
```

**风险规则示例**:
```go
var AddressValidationRule = RiskRule{
    Name:   "地址验证",
    Weight: 0.15,
    Calculate: func(order *Order, profile *UserProfile) float64 {
        if order.SenderAddress == order.ReceiverAddress {
            return 0  // 高风险
        }
        if !ValidateAddress(order.SenderAddress) {
            return 0  // 高风险
        }
        return 15  // 正常
    },
}
```

#### 4.4 风控API设计

**风险评估**: `POST /api/v1/risk/evaluate`

**请求参数**:
```go
type EvaluateRiskRequest struct {
    OrderId int64 `json:"order_id"`
}
```

**响应结果**:
```go
type EvaluateRiskResponse struct {
    RiskScore   int      `json:"risk_score"`
    RiskLevel   string   `json:"risk_level"`
    RiskReasons []string `json:"risk_reasons"`
    Action      string   `json:"action"`
}
```

### 5. 费用计算引擎优化模块

#### 5.1 组件设计

```
internal/
├── biz/
│   ├── pricing.go             # 费用计算逻辑（已存在）
│   └── pricing_rule.go        # 计费规则引擎
├── data/
│   └── pricing_cache.go       # 费用缓存
└── service/
    └── pricing_service.go     # 费用计算服务
```

#### 5.2 规则引擎设计

**规则定义**:
```go
type PricingRule struct {
    Name      string
    Condition func(*Order) bool
    Calculate func(*Order) float64
}

type PricingEngine struct {
    BaseCalculator func(*Order) float64
    Rules          []PricingRule
}

func (pe *PricingEngine) Calculate(order *Order) float64 {
    totalFee := pe.BaseCalculator(order)
    for _, rule := range pe.Rules {
        if rule.Condition(order) {
            totalFee += rule.Calculate(order)
        }
    }
    return totalFee
}
```

**规则示例**:
```go
var InsuranceRule = PricingRule{
    Name: "保价费",
    Condition: func(order *Order) bool {
        return order.InsuredValue > 0
    },
    Calculate: func(order *Order) float64 {
        fee := order.InsuredValue * 0.005
        if fee < 2 {
            return 2
        }
        return fee
    },
}
```

#### 5.3 缓存策略

**缓存Key设计**:
```
pricing:cache:{origin}:{destination}:{weight}:{service_type}
```

**缓存内容**:
```go
type PricingCache struct {
    BaseFee      float64 `json:"base_fee"`
    InsuranceFee float64 `json:"insurance_fee"`
    TotalFee     float64 `json:"total_fee"`
    ValidUntil   string  `json:"valid_until"`
}
```

**缓存时长**: 1小时

**缓存更新**: 规则变更时清除相关缓存

#### 5.4 费用计算API设计

**计算费用**: `POST /api/v1/pricing/calculate`

**请求参数**:
```go
type CalculatePricingRequest struct {
    Origin      string  `json:"origin"`
    Destination string  `json:"destination"`
    Weight      float64 `json:"weight"`
    ServiceType string  `json:"service_type"`
    InsuredValue float64 `json:"insured_value"`
}
```

**响应结果**:
```go
type CalculatePricingResponse struct {
    BaseFee      float64            `json:"base_fee"`
    ExtraFees    map[string]float64 `json:"extra_fees"`
    TotalFee     float64            `json:"total_fee"`
    Cached       bool               `json:"cached"`
}
```

## 数据流设计

### 订单搜索数据流

```
用户查询 → API网关 → 搜索服务 → Elasticsearch → 返回结果
                                ↓
                            权限过滤 + 脱敏
```

### 订单导出数据流

```
用户请求 → 创建任务 → Redis Stream → Worker消费
                                      ↓
                                  分批查询ES/MySQL
                                      ↓
                                  生成Excel
                                      ↓
                                  上传OSS
                                      ↓
                                  通知用户
```

### 补偿流程数据流

```
异常上报 → 规则引擎计算 → 判断是否需要审批
                            ↓
                    需要审批 → 人工审核 → 审批通过 → 打款
                            ↓
                    自动审批 → 打款 → 通知用户
```

### 风控评估数据流

```
订单创建 → 风控评估 → 规则引擎计算风险分
                        ↓
                  低风险 → 正常流转
                  中风险 → 人工审核
                  高风险 → 拦截订单
```

## 错误处理

### 错误码设计

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 40001 | 参数错误 | 400 |
| 40002 | 订单不存在 | 404 |
| 40003 | 权限不足 | 403 |
| 50001 | ES查询失败 | 500 |
| 50002 | 导出任务创建失败 | 500 |
| 50003 | 补偿计算失败 | 500 |
| 50004 | 风控评估失败 | 500 |
| 50005 | 费用计算失败 | 500 |

### 降级策略

1. **ES查询降级**: ES不可用时降级到MySQL查询
2. **缓存降级**: Redis不可用时直接计算
3. **异步任务降级**: 队列满时拒绝新任务
4. **风控降级**: 风控服务异常时允许订单通过但记录日志

## 测试策略

### 单元测试

- 规则引擎逻辑测试
- 数据转换测试
- 缓存逻辑测试

### 集成测试

- ES数据同步测试
- 导出任务端到端测试
- 补偿流程测试
- 风控评估测试

### 性能测试

- ES查询性能测试（目标: <500ms）
- 导出速度测试（目标: >1000条/秒）
- 风控评估性能测试（目标: <100ms）
- 费用计算性能测试（目标: <20ms）

### 压力测试

- 并发查询压力测试
- 大批量导出压力测试
- 高并发风控评估测试

## 部署方案

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Elasticsearch 8.0+
- 阿里云OSS

### 部署步骤

1. 部署Elasticsearch集群
2. 配置Canal同步服务
3. 部署后端服务
4. 启动Worker进程
5. 配置监控告警

### 监控指标

- ES查询QPS和响应时间
- 导出任务成功率和处理速度
- 补偿审批通过率和处理时长
- 风控拦截率和误杀率
- 费用计算缓存命中率

## 安全考虑

### 数据安全

- 敏感信息脱敏展示
- 导出文件加密存储
- 访问日志记录

### 权限控制

- 基于角色的访问控制（RBAC）
- 订单查询权限分级
- 补偿审批权限管理
- 风控规则修改权限

### 防护措施

- API限流
- 防重复提交
- SQL注入防护
- XSS防护
