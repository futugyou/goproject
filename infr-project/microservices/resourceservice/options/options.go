package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Options struct {
	DBName            string
	MongoDBURL        string
	QstashToken       string
	QstashDestination string
}

func New() *Options {
	opts := &Options{
		DBName:            os.Getenv("db_name"),
		MongoDBURL:        os.Getenv("mongodb_url"),
		QstashToken:       os.Getenv("QSTASH_TOKEN"),
		QstashDestination: os.Getenv("QSTASH_DESTINATION"),
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
