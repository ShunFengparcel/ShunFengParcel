package config

import (
	"time"

	"gorm.io/gorm"
)

type Reconciliation struct {
	gorm.Model
	TaskName           string    `gorm:"type:varchar(100);NOT NULL;COMMENT: 任务名称,"`
	ReconciliationNo   string    `gorm:"type:varchar(100);not null;comment:对账编号"`
	ReconciliationDate time.Time `gorm:"type:datetime;not null;index;comment:对账日期"`
	Amount             float64   `gorm:"type:decimal(10,2);not null;comment:对账金额"`
	ActualAmount       float64   `gorm:"type:decimal(10,2);not null;comment:实际金额"`
	Count              int64     `gorm:"type:int(11);not null;comment:完成数量"`
	Proportion         float64   `gorm:"type:decimal(10,2);not null;comment:涨幅比例"`
	HandlerStatus      string    `gorm:"type:enum('待处理','处理中','已完成');not null;default:'待处理';comment:处理状态 '待处理','处理中','已完成'"`
}

func (r Reconciliation) TableName() string {
	return "reconciliation"
}

func (r *Reconciliation) Created(DB *gorm.DB) error {
	return DB.Create(&r).Error
}

func (r *Reconciliation) ReconciliationItemList(DB *gorm.DB, playtime, endtime string) (list []Reconciliation, err error) {
	query := DB.Model(&Reconciliation{})

	if playtime != "" && endtime != "" {
		query = query.Where("created_at >= ? and created_at <= ?", playtime, endtime)
	}

	err = query.Order("created_at desc").Find(&list).Error
	return list, err
}
