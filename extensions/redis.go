package extensions

import (
	"context"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

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

var HashAllLueScript2 = redis.NewScript(`
	local results = {}
	local missing_keys = {}

	local updated_keys = redis.call("SMEMBERS", ARGV[1])

	for _, key in ipairs(updated_keys) do
		if redis.call("EXISTS", key) == 1 then
			redis.call("DEL", key)
		else
			table.insert(missing_keys, key)
		end
	end

	redis.call("DEL", ARGV[1])

	local scan_result = redis.call("SCAN", 0, "MATCH", ARGV[2], "COUNT", ARGV[3])
	for _, key in ipairs(scan_result[2]) do
		local hash = redis.call("HGETALL", key)
		local result = {}
		for i = 1, #hash, 2 do
			result[hash[i]] = hash[i+1]
		end
		table.insert(results, result)
	end

	return {results, missing_keys}
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

// Updated_key is a Set that stores expired keys. eg. create/update/delete key
// missingKeys can be used to supplement missing data from the database
func RedisListHashWithLua2[T any](ctx context.Context, client *redis.Client, pattern string, count int64, updated_key string) ([]T, []string, error) {
	pattern = pattern + "*"

	res, err := HashAllLueScript2.Run(ctx, client, nil, updated_key, pattern, count).Result()
	if err != nil {
		return nil, nil, err
	}

	rawResults := res.([]interface{})[0].([]interface{})
	rawMissingKeys := res.([]interface{})[1].([]interface{})

	var datas []T
	for _, item := range rawResults {
		data := item.(map[string]interface{})

		var rv T
		err := mapstructure.Decode(data, &rv)
		if err != nil {
			return nil, nil, err
		}

		datas = append(datas, rv)
	}

	var missingKeys []string
	for _, key := range rawMissingKeys {
		missingKeys = append(missingKeys, key.(string))
	}

	return datas, missingKeys, nil
}

func GetLock(ctx context.Context, client *redis.Client, lockKey string, lockValue string, lockTTL time.Duration) error {
	ok, err := client.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to acquire lock")
	}
	return nil
}

var releaseLockScript = redis.NewScript(`
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
`)

func ReleaseLock(ctx context.Context, client *redis.Client, lockKey string, lockValue string) (int64, error) {
	result, err := releaseLockScript.Run(ctx, client, []string{lockKey}, lockValue).Result()
	if err != nil {
		return 0, err
	}

	if success, ok := result.(int64); ok {
		return success, nil
	} else {
		return 0, errors.New("redis data type error, check script and code")
	}
}
