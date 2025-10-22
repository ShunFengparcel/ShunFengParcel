package data

import (
	"ShunFengParcel/internal/conf"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 从配置生成 Redis 客户端，交给 Wire 注入
func NewRedisClient(c *conf.Data) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:     "14.103.136.136:6379",
		Password: "b914572bdc69c888f90ce7bafa9cbb7e", // no password set
		DB:       0,
	})
}
