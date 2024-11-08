package extensions

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	opt.MaxRetries = 3
	opt.DialTimeout = 10 * time.Second
	opt.ReadTimeout = -1
	opt.WriteTimeout = -1
	opt.DB = 0

	return redis.NewClient(opt), nil
}
