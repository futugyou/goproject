package services

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Rbd *redis.Client
var ctx = context.Background()

func init() {
	redisAddress := os.Getenv("redisAddress")
	redisPassword := os.Getenv("redisAddress")
	redisDBString := os.Getenv("redisAddress")
	redisDB, _ := strconv.Atoi(redisDBString)
	Rbd = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDB,
	})
}
