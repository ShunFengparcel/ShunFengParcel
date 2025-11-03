package config

import (
	"time"

	"gorm.io/gorm"
)

type SfPayments struct {
	Id        int64     `gorm:"type:int(11);comment:支付id"`                                           // 支付ID
	OrderId   int64     `gorm:"type:int(11);comment:支付id"`                                           // 订单ID
	PayNo     string    `gorm:"type:varchar(50);not null;comment:支付单号"`                              // 支付单号
	Channel   string    `gorm:"type:enum('wechat','alipay','cash','monthly');not null;comment:支付渠道"` // 渠道
	Amount    float64   `gorm:"type:decimal(10,2);not null;comment:支付金额"`                            // 金额
	Currency  string    `gorm:"type:varchar(3);default:'CNY';not null"`
	Status    string    `gorm:"type:enum('success','failed','processing');not null;comment:支付状态"` // 状态
	PaidAt    time.Time `gorm:"type:datetime;not null;comment;支付成功时间"`                            // 支付成功时间
	ThirdTxId string    `gorm:"type:varchar(50);not null;comment:第三方流水号"`                         // 第三方流水号
	CreatedAt time.Time `gorm:"type:datetime;not null;comment:创建时间"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:修改时间"`
	DeleteAt  time.Time `gorm:"type:datetime;comment:删除时间"` // 删除时间
}

func (SfPayments) TableName() string {
	return "sf_payments"
}

func (p SfPayments) Created(DB *gorm.DB) error {
	return DB.Create(&p).Error
}

func (p SfPayments) Updated(DB *gorm.DB) error {
	return DB.Updates(&p).Error
}

func (p SfPayments) FIndByList(DB *gorm.DB, status string) (List []*SfPayments, err error) {
	a := DB.Model(p)
	if status != "" {
		a = a.Where("handler_status = ?", status)
	}
	err = a.Order("created_at desc").Find(&List).Error

	return
}
