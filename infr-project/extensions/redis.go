package extensions

import (
	"context"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
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

func RedisScanHashAll[T any](ctx context.Context, client *redis.Client, prefix string, count int64) ([]T, error) {
	var datas []T
	iter := client.Scan(ctx, 0, prefix+"*", count).Iterator()

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

var HashAllLueScript = redis.NewScript(`
	local keys = redis.call("SCAN", 0, "MATCH", ARGV[1], "COUNT", ARGV[2])
	local results = {}
	for _, key in ipairs(keys[2]) do
		local hash = redis.call("HGETALL", key)
		local result = {}
		for i = 1, #hash, 2 do
			result[hash[i]] = hash[i+1]
		end
		table.insert(results, result)
	end
	return results
`)

func RedisListHashWithLua[T any](ctx context.Context, client *redis.Client, pattern string, count int64) ([]T, error) {
	pattern = pattern + "*"

	res, err := HashAllLueScript.Run(ctx, client, nil, pattern, count).Result()
	if err != nil {
		return nil, err
	}

	var datas []T
	for _, item := range res.([]interface{}) {
		data := item.(map[string]interface{})

		var rv T
		err := mapstructure.Decode(data, &rv)
		if err != nil {
			return nil, err
		}

		datas = append(datas, rv)
	}

	return datas, nil
}
