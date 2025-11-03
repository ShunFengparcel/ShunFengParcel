package config

import "gorm.io/gorm"

type TransactionMonitor struct {
	gorm.Model
	MonitorNo         string  `json:"type:varchar(50);not null;comment:监控编号"`   // 监控编号
	OrderNo           string  `json:"type:varchar(50);not null;comment:订单编号"`   // 订单号
	UserId            int64   `json:"type:int(11);not null;comment:用户id"`       // 用户ID
	TransactionAmount float64 `json:"type:decimal(10,2);not null;comment:交易金额"` // 交易金额
	MonitorType       string  `json:"type:varchar(50);not null;comment:监控类型"`   // 监控类型
	RiskLevel         string  `json:"type:enum('low','medium','high','critical');low;not null;comment:监控等级"`
	RiskScore         int64   `json:"type:int(11);not null;comment:风险分数"` // 风险分数
	MonitorResult     string  `json:"type:enum('pass','review','block');default:pass;not null;comment:监控结果"`
	HandleStatus      string  `json:"enum('pending','processing','handled');default:pending;not null;comment:处理类型"`
	HandleResult      string  `json:"type:varchar(100);comment:处理结果"`
}

func (m *TransactionMonitor) TableName() string {
	return "transaction_monitor"
}

func (m *TransactionMonitor) FindByID(DB *gorm.DB, id int64) (err error) {
	return DB.Where("id = ?", id).First(m).Error
}

func (m *TransactionMonitor) Created(DB *gorm.DB) error {
	return DB.Create(&m).Error
}

func (m *TransactionMonitor) Updated(DB *gorm.DB, monitorNo string) error {
	return DB.Model(m).Where("monitor_no = ?", monitorNo).Updates(&m).Error
}

func (m *TransactionMonitor) Deleted(DB *gorm.DB, ID int) error {
	return DB.Where("id = ?", ID).Delete(&m).Error
}
