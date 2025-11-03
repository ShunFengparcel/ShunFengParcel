package config

import (
	"time"

	"gorm.io/gorm"
)

type SfOrders struct {
	Id                   int64     `gorm:"column:id;type:bigint;comment:订单ID;primaryKey;not null;" json:"id"`                                                                                                                                     // 订单ID
	OrderNo              string    `gorm:"column:order_no;type:varchar(20);comment:订单号;not null;" json:"order_no"`                                                                                                                                // 订单号
	UserId               int64     `gorm:"column:user_id;type:bigint;comment:用户ID;not null;" json:"user_id"`                                                                                                                                      // 用户ID
	CourierId            int64     `gorm:"column:courier_id;type:bigint;comment:快递员ID;default:NULL;" json:"courier_id"`                                                                                                                           // 快递员ID
	SenderName           string    `gorm:"column:sender_name;type:varchar(50);comment:寄件人;not null;" json:"sender_name"`                                                                                                                          // 寄件人
	SenderPhone          string    `gorm:"column:sender_phone;type:varchar(20);comment:寄件电话;not null;" json:"sender_phone"`                                                                                                                       // 寄件电话
	SenderAddress        string    `gorm:"column:sender_address;type:json;comment:寄件地址;not null;" json:"sender_address"`                                                                                                                          // 寄件地址
	ReceiverName         string    `gorm:"column:receiver_name;type:varchar(50);comment:收件人;not null;" json:"receiver_name"`                                                                                                                      // 收件人
	ReceiverPhone        string    `gorm:"column:receiver_phone;type:varchar(20);comment:收件电话;not null;" json:"receiver_phone"`                                                                                                                   // 收件电话
	ReceiverAddress      string    `gorm:"column:receiver_address;type:json;comment:收件地址;not null;" json:"receiver_address"`                                                                                                                      // 收件地址
	ServiceType          string    `gorm:"column:service_type;type:enum('same_day', 'next_day', 'standard', 'economy');comment:服务类型;not null;" json:"service_type"`                                                                               // 服务类型
	ProductType          string    `gorm:"column:product_type;type:varchar(50);comment:物品类型;default:NULL;" json:"product_type"`                                                                                                                   // 物品类型
	IsInsured            int8      `gorm:"column:is_insured;type:tinyint(1);comment:是否保价;default:0;" json:"is_insured"`                                                                                                                           // 是否保价
	InsuredValue         float64   `gorm:"column:insured_value;type:decimal(10, 2);comment:保价金额;default:NULL;" json:"insured_value"`                                                                                                              // 保价金额
	EstimatedWeight      float64   `gorm:"column:estimated_weight;type:decimal(8, 3);comment:预估重量kg;default:NULL;" json:"estimated_weight"`                                                                                                       // 预估重量kg
	ActualWeight         float64   `gorm:"column:actual_weight;type:decimal(8, 3);comment:实际重量kg;default:NULL;" json:"actual_weight"`                                                                                                             // 实际重量kg
	VolumeWeight         float64   `gorm:"column:volume_weight;type:decimal(8, 3);comment:体积重量kg;default:NULL;" json:"volume_weight"`                                                                                                             // 体积重量kg
	ChargeWeight         float64   `gorm:"column:charge_weight;type:decimal(8, 3);comment:计费重量;" json:"charge_weight"`                                                                                                                            // 计费重量
	EstimatedFee         float64   `gorm:"column:estimated_fee;type:decimal(10, 2);comment:预估费用;default:NULL;" json:"estimated_fee"`                                                                                                              // 预估费用
	ActualFee            float64   `gorm:"column:actual_fee;type:decimal(10, 2);comment:实际费用;default:NULL;" json:"actual_fee"`                                                                                                                    // 实际费用
	PaymentMethod        string    `gorm:"column:payment_method;type:enum('sender_pay', 'receiver_pay', 'monthly');comment:支付方式;not null;" json:"payment_method"`                                                                                 // 支付方式
	PaymentStatus        string    `gorm:"column:payment_status;type:enum('pending', 'paid', 'failed', 'refunded');comment:支付状态;not null;default:pending;" json:"payment_status"`                                                                 // 支付状态
	ExpectedPickupTime   time.Time `gorm:"column:expected_pickup_time;type:datetime;comment:期望上门时间;default:NULL;" json:"expected_pickup_time"`                                                                                                    // 期望上门时间
	ActualPickupTime     time.Time `gorm:"column:actual_pickup_time;type:datetime;comment:实际上门时间;default:NULL;" json:"actual_pickup_time"`                                                                                                        // 实际上门时间
	ExpectedDeliveryTime time.Time `gorm:"column:expected_delivery_time;type:datetime;comment:期望送达时间;default:NULL;" json:"expected_delivery_time"`                                                                                                // 期望送达时间
	ActualDeliveryTime   time.Time `gorm:"column:actual_delivery_time;type:datetime;comment:实际送达时间;default:NULL;" json:"actual_delivery_time"`                                                                                                    // 实际送达时间
	OrderStatus          string    `gorm:"column:order_status;type:enum('pending', 'accepted', 'picked_up', 'in_transit', 'out_for_delivery', 'delivered', 'cancelled', 'exception');comment:订单状态;not null;default:pending;" json:"order_status"` // 订单状态
	CreatedAt            time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeleteAt             time.Time `gorm:"column:delete_at;type:datetime;comment:删除时间;default:NULL;" json:"delete_at"` // 删除时间
}

func (o *SfOrders) UpdateOrderStatus(DB *gorm.DB, orderno string) error {
	return DB.Model(o).Where("order_no = ?", orderno).Updates(&o).Error
}

func (o *SfOrders) FIndByOrderSn(DB *gorm.DB, sn string) error {
	return DB.Model(o).Where("order_no = ?", sn).Limit(1).Find(&o).Error
}

func (o *SfOrders) FIndByList(DB *gorm.DB, OrderStatus string) (List []*SfOrders, err error) {
	query := DB.Model(o)
	if OrderStatus != "" {
		query = query.Where("order_status = ?", OrderStatus)
	}

	err = query.Order("created_at desc").Find(&List).Error
	return
}
