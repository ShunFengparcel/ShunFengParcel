package model

import (
	"time"
)

type SfCourierTasks struct {
	Id         int64     `gorm:"column:id;type:bigint;comment:任务ID;primaryKey;not null;" json:"id"`                                                              // 任务ID
	OrderId    int64     `gorm:"column:order_id;type:bigint;comment:订单id;default:NULL;" json:"order_id"`                                                         // 订单id
	CourierId  int64     `gorm:"column:courier_id;type:bigint;comment:快递员ID;not null;" json:"courier_id"`                                                        // 快递员ID
	Priority   int8      `gorm:"column:priority;type:tinyint;comment:优先级 0最高 1中等 2最低;default:1;" json:"priority"`                                                // 优先级 0最高 1中等 2最低
	TaskStatus string    `gorm:"column:task_status;type:enum('pending', 'accepted', 'completed', 'cancelled');comment:任务状态;default:pending;" json:"task_status"` // 任务状态
	TaskType   string    `gorm:"column:task_type;type:enum('pickup', 'delivery');comment:任务类型;not null;" json:"task_type"`                                       // 任务类型
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeleteAt   time.Time `gorm:"column:delete_at;type:datetime;comment:删除时间;default:NULL;" json:"delete_at"` // 删除时间
}

func (SfCourierTasks) TableName() string {
	return "sf_courier_tasks"
}
