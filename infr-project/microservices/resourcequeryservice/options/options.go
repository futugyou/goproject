package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Options struct {
	RedisURL        string
	QueryDBName     string
	QueryMongoDBURL string
}

func New() *Options {
	opts := &Options{
		RedisURL:        os.Getenv("REDIS_URL"),
		QueryDBName:     os.Getenv("query_db_name"),
		QueryMongoDBURL: os.Getenv("query_mongodb_url"),
	}

	return opts
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
