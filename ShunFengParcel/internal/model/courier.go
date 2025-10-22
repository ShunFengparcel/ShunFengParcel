package model

import (
	"time"
)

type SfCouriers struct {
	Id               int64     `gorm:"column:id;type:bigint;comment:快递员ID;primaryKey;not null;" json:"id"`                                 // 快递员ID
	EmployeeNo       string    `gorm:"column:employee_no;type:varchar(20);comment:工号;not null;" json:"employee_no"`                        // 工号
	RealName         string    `gorm:"column:real_name;type:varchar(50);comment:姓名;not null;" json:"real_name"`                            // 姓名
	Phone            string    `gorm:"column:phone;type:varchar(20);comment:电话;not null;" json:"phone"`                                    // 电话
	PasswordHash     string    `gorm:"column:password_hash;type:varchar(255);comment:密码哈希;not null;" json:"password_hash"`                 // 密码哈希
	RegionCode       string    `gorm:"column:region_code;type:varchar(20);comment:所属区域编码;not null;" json:"region_code"`                    // 所属区域编码
	Status           string    `gorm:"column:status;type:enum('active', 'inactive', 'resigned');comment:状态;default:active;" json:"status"` // 状态
	PerformanceScore float64   `gorm:"column:performance_score;type:decimal(5, 4);comment:综合绩效分;default:0.0000;" json:"performance_score"` // 综合绩效分
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeleteAt         time.Time `gorm:"column:delete_at;type:datetime;comment:删除时间;default:NULL;" json:"delete_at"` // 删除时间
}

func (SfCouriers) TableName() string {
	return "sf_couriers"
}
