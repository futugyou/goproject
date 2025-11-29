package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Options struct {
	DBName           string
	MongoDBURL       string
	SendgridApiKey   string
	EmailVerifyUrl   string
	EmailFromAddress string
	EmailFromName    string
}

func New() *Options {
	opts := &Options{
		DBName:           os.Getenv("db_name"),
		MongoDBURL:       os.Getenv("mongodb_url"),
		SendgridApiKey:   os.Getenv("SENDGRID_API_KEY"),
		EmailVerifyUrl:   os.Getenv("EMAIL_VERIFY_URL"),
		EmailFromAddress: os.Getenv("EMAIL_FROM_ADDRESS"),
		EmailFromName:    os.Getenv("EMAIL_FROM_NAME"),
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
