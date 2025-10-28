package data

import (
	"ShunFengParcel/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo, NewOrderRepo)

// Data .
type Data struct {
	DB  *gorm.DB
	Rdb *redis.Client

	// TODO wrapped database client
}

var (
	Ctx = context.Background()
)

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password, // no password set
		DB:       0,                // use default DB
	})

	// 测试 Redis 连接（如果失败只记录日志，不中断启动）
	err := rdb.Ping(Ctx).Err()
	if err != nil {
		log.NewHelper(logger).Warnf("Redis connection failed: %v, continuing without Redis", err)
		// 不 panic，继续运行
	} else {
		log.NewHelper(logger).Info("Redis connected successfully")
	}

	// 配置 GORM
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, nil, err
	}

	// 获取底层的 sql.DB 对象进行连接池配置
	_, err = db.DB()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close() // 关闭 Mysql 连接
		}
		rdb.Close() // 关闭 Redis 连接
	}

	return &Data{DB: db, Rdb: rdb}, cleanup, nil
}
