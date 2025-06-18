package extensions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
		table.insert(results, cjson.encode(result))  -- 👈 Key: Convert to JSON string
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
		table.insert(results, cjson.encode(result)) -- 👈 Convert to JSON string
	end

	return {results, missing_keys}
`)

func RedisListHashWithLua[T any](ctx context.Context, client *redis.Client, pattern string, count int64) ([]T, error) {
	pattern = pattern + "*"

	res, err := HashAllLueScript.Run(ctx, client, nil, pattern, count).Result()
	if err != nil {
		return nil, err
	}
	rawItems, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected Redis result type: %T", res)
	}

	var datas []T
	for _, item := range rawItems {
		jsonStr, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected item type: %T", item)
		}

		var rv T
		err := json.Unmarshal([]byte(jsonStr), &rv)
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

	// First assert that res is []interface{}
	resList, ok := res.([]interface{})
	if !ok || len(resList) != 2 {
		return nil, nil, fmt.Errorf("unexpected result structure from Lua script: %T", res)
	}

	// Get the result part
	rawResults, ok := resList[0].([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("unexpected result[0] type: %T", resList[0])
	}

	rawMissingKeys, ok := resList[1].([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("unexpected result[1] type: %T", resList[1])
	}

	var datas []T
	for _, item := range rawResults {
		jsonStr, ok := item.(string)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected item type: %T", item)
		}

		var rv T
		if err := json.Unmarshal([]byte(jsonStr), &rv); err != nil {
			return nil, nil, err
		}
		datas = append(datas, rv)
	}

	var missingKeys []string
	for i, key := range rawMissingKeys {
		strKey, ok := key.(string)
		if !ok {
			return nil, nil, fmt.Errorf("unexpected key type at index %d: %T", i, key)
		}
		missingKeys = append(missingKeys, strKey)
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
