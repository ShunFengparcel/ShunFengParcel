package config

import (
	"ShunFengParcel/api/helloworld/payment"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Reconciliation struct {
	gorm.Model
	TaskName           string    `gorm:"type:varchar(100);NOT NULL;COMMENT: 任务名称,"`
	ReconciliationNo   string    `gorm:"type:varchar(100);not null;comment:对账编号"`
	ReconciliationDate time.Time `gorm:"type:datetime;not null;index;comment:对账日期"`
	Amount             float64   `gorm:"type:decimal(10,2);not null;comment:对账金额"`
	ActualAmount       float64   `gorm:"type:decimal(10,2);not null;comment:对账金额"`
	Proportion         float64   `gorm:"type:decimal(10,2);not null;comment:涨幅比例"`
	HandlerStatus      string    `gorm:"type:enum('待处理','处理中','已完成');not null;default:'待处理';comment:处理状态 '待处理','处理中','已完成'"`
}

func (r Reconciliation) TableName() string {
	return "reconciliation"
}

func (r *Reconciliation) Created(DB *gorm.DB) error {
	return DB.Create(&r).Error
}

func (r *Reconciliation) ReconciliationItemList(DB *gorm.DB, playtime, endtime string) (list []*payment.ReconciliationItem, err error) {
	var reconciliations []Reconciliation
	query := DB.Model(&Reconciliation{})

	if playtime != "" && endtime != "" {
		query = query.Where("created_at >= ? and created_at <= ?", playtime, endtime)
	}

	err = query.Order("created_at desc").Find(&reconciliations).Error
	if err != nil {
		return nil, err
	}

	// 将 Reconciliation 转换为 payment.ReconciliationItem
	for _, rec := range reconciliations {
		item := &payment.ReconciliationItem{
			TaskName:           rec.TaskName,
			ReconciliationNo:   rec.ReconciliationNo,
			ReconciliationDate: rec.ReconciliationDate.Format("2006-01-02 15:04:05"),
			Amount:             float32(rec.Amount),
			ActualAmount:       float32(rec.ActualAmount),
			Proportion:         fmt.Sprintf("%.2f", rec.Proportion),
			HandlerStatus:      rec.HandlerStatus,
		}
		list = append(list, item)
	}

	return list, nil
}
