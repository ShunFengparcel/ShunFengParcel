package model

import "time"

// SfOrderAuditLogs 订单操作审计日志
type SfOrderAuditLogs struct {
    Id             int64     `gorm:"column:id;type:bigint;primaryKey;not null" json:"id"`
    OrderId        int64     `gorm:"column:order_id;type:bigint;index;not null" json:"order_id"`
    ActionType     string    `gorm:"column:action_type;type:varchar(50);not null" json:"action_type"` // cancel / reassign / exception
    OperatorId     int64     `gorm:"column:operator_id;type:bigint;default:NULL" json:"operator_id"`
    OperatorRole   string    `gorm:"column:operator_role;type:varchar(50);default:NULL" json:"operator_role"` // courier/system/admin
    BeforeStatus   string    `gorm:"column:before_status;type:varchar(50);default:NULL" json:"before_status"`
    AfterStatus    string    `gorm:"column:after_status;type:varchar(50);default:NULL" json:"after_status"`
    ReasonType     string    `gorm:"column:reason_type;type:varchar(50);default:NULL" json:"reason_type"`
    ReasonText     string    `gorm:"column:reason_text;type:varchar(255);default:NULL" json:"reason_text"`
    IdempotencyKey string    `gorm:"column:idempotency_key;type:varchar(100);default:NULL" json:"idempotency_key"`
    Metadata       string    `gorm:"column:metadata;type:text;default:NULL" json:"metadata"`
    CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (SfOrderAuditLogs) TableName() string { return "sf_order_audit_logs" }