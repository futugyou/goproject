package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Options struct {
	DBName     string
	MongoDBURL string
}

func New() *Options {
	opts := &Options{
		DBName:     os.Getenv("db_name"),
		MongoDBURL: os.Getenv("mongodb_url"),
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
