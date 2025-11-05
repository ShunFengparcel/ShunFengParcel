package model

import "time"

// SfOrderReassignments 改派链路记录
type SfOrderReassignments struct {
	Id              int64     `gorm:"column:id;type:bigint;primaryKey;not null" json:"id"`
	OrderId         int64     `gorm:"column:order_id;type:bigint;index;not null" json:"order_id"`
	FromCourierId   int64     `gorm:"column:from_courier_id;type:bigint;default:NULL" json:"from_courier_id"`
	ToCourierId     int64     `gorm:"column:to_courier_id;type:bigint;not null" json:"to_courier_id"`
	OperatorId      int64     `gorm:"column:operator_id;type:bigint;default:NULL" json:"operator_id"`
	ReasonText      string    `gorm:"column:reason_text;type:varchar(255);default:NULL" json:"reason_text"`
	DispatchVersion int64     `gorm:"column:dispatch_version;type:bigint;default:0" json:"dispatch_version"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (SfOrderReassignments) TableName() string { return "sf_order_reassignments" }
