package services

import (
	"context"

	"github.com/beego/beego/v2/core/config"

	"github.com/redis/go-redis/v9"
)

var Rbd *redis.Client
var ctx = context.Background()

func init() {
	redisAddress, _ := config.String("redisAddress")
	redisPassword, _ := config.String("redisAddress")
	redisDB, _ := config.Int("redisAddress")
	Rbd = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDB,
	})
}
