package inits

import (
	"ShunFengParcel/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMysql() {
	var err error
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:mysql_8tDfmn@tcp(14.103.136.136:3306)/shunfeng?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&config.Reconciliation{}, &config.TransactionMonitor{}, &config.SfPayments{})
}
