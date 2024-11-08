package extensions

import (
	"context"
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

func RedisScanHashAll[T any](ctx context.Context, prefix string, client *redis.Client) ([]T, error) {
	var datas []T
	iter := client.Scan(ctx, 0, prefix+"*", 100).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		var rv T
		err := client.HGetAll(ctx, key).Scan(&rv)
		if err != nil {
			return nil, err
		}
		datas = append(datas, rv)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return datas, nil
}
