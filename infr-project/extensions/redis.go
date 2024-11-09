package extensions

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func RedisClient(url string) (*redis.Client, error) {
	var err error
	once.Do(func() {
		opt, err := redis.ParseURL(url)
		if err != nil {
			return
		}

		opt.MaxRetries = 3
		opt.DialTimeout = 10 * time.Second
		opt.ReadTimeout = -1
		opt.WriteTimeout = -1
		opt.DB = 0

		redisClient = redis.NewClient(opt)
	})

	return redisClient, err
}
