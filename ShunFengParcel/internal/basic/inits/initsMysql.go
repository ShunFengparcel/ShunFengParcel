package inits

import (
	"ShunFengParcel/internal/basic/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func InitsMysql() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:mysql_8tDfmn@tcp(14.103.136.136:3306)/shunfeng?charset=utf8mb4&parseTime=True&loc=Local"
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		log.Println("数据库连接成功")
	}
}
