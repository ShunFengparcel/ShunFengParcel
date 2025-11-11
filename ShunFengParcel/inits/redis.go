package inits

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "14.103.175.138:6379",
		Password: "3c907a418fdb4cf3112ed3b118669ed4", // no password set
		DB:       0,                                  // use default DB
	})

	err := RDB.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("redis init success")
	}
}
