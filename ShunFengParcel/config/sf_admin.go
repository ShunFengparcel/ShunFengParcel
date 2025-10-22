package config

import "time"

type SysAdmin struct {
	Id            int64     `gorm:"column:id;type:bigint;comment:主键ID;primaryKey;not null;" json:"id"`                                                                   // 主键ID
	Username      string    `gorm:"column:username;type:varchar(50);comment:登录用户名（唯一）;not null;" json:"username"`                                                        // 登录用户名（唯一）
	Password      string    `gorm:"column:password;type:varchar(255);comment:登录密码（加密存储，如bcrypt/MD5）;not null;" json:"password"`                                          // 登录密码（加密存储，如bcrypt/MD5）
	Nickname      string    `gorm:"column:nickname;type:varchar(100);comment:管理员昵称;default:NULL;" json:"nickname"`                                                       // 管理员昵称
	RoleId        int64     `gorm:"column:role_id;type:bigint;comment:关联角色表ID（用于权限控制）;not null;" json:"role_id"`                                                         // 关联角色表ID（用于权限控制）
	Status        int8      `gorm:"column:status;type:tinyint(1);comment:状态：1-启用，0-禁用;not null;default:1;" json:"status"`                                                // 状态：1-启用，0-禁用
	LastLoginTime time.Time `gorm:"column:last_login_time;type:datetime;comment:最后登录时间;default:NULL;" json:"last_login_time"`                                            // 最后登录时间
	Phone         string    `gorm:"column:phone;type:varchar(20);comment:联系电话;default:NULL;" json:"phone"`                                                               // 联系电话
	Email         string    `gorm:"column:email;type:varchar(100);comment:邮箱;default:NULL;" json:"email"`                                                                // 邮箱
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;comment:创建时间;not null;default:CURRENT_TIMESTAMP;" json:"create_time"`                                // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;comment:更新时间（低版本MySQL不支持ON UPDATE，可通过应用层更新）;not null;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间（低版本MySQL不支持ON UPDATE，可通过应用层更新）
}

func (SysAdmin) TableName() string {
	return "sys_admin"
}
