package inits

import (
	"ShunFengParcel/internal/basic/config"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitsRedis() {
	config.RDB = redis.NewClient(&redis.Options{
		Addr:     "14.103.136.136:6379",
		Password: "b914572bdc69c888f90ce7bafa9cbb7e", // no password set
		DB:       0,                                  // use default DB
	})

	err = config.RDB.Set(Ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	log.Println("init redis success")
}
