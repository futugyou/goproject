package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var GlobalOptions *Options

type Options struct {
	DBName     string
	MongoDBURL string

	QstashToken       string
	QstashDestination string

	AwsAccessKeyID     string
	AwsSecretAccessKey string
	AwsRegion          string

	HcpClientID       string
	HcpClientSecret   string
	HcpAppName        string
	HcpOrganizationID string
	HcpProjectID      string

	AzureClientID     string
	AzureTenantID     string
	AzureClientSecret string
	AzureVaultURL     string

	EncryptKey string
}

func init() {
	if GlobalOptions == nil {
		GlobalOptions = New()
	}
}

func New() *Options {
	opts := &Options{
		DBName:     os.Getenv("db_name"),
		MongoDBURL: os.Getenv("mongodb_url"),

		QstashToken:       os.Getenv("QSTASH_TOKEN"),
		QstashDestination: os.Getenv("QSTASH_DESTINATION"),

		AwsAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsRegion:          os.Getenv("AWS_REGION"),

		HcpClientID:       os.Getenv("HCP_CLIENT_ID"),
		HcpClientSecret:   os.Getenv("HCP_CLIENT_SECRET"),
		HcpAppName:        os.Getenv("HCP_APP_NAME"),
		HcpOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		HcpProjectID:      os.Getenv("HCP_PROJECT_ID"),

		AzureClientID:     os.Getenv("AZURE_CLIENT_ID"),
		AzureTenantID:     os.Getenv("AZURE_TENANT_ID"),
		AzureClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		AzureVaultURL:     os.Getenv("AZURE_VAULT_URL"),

		EncryptKey: os.Getenv("Encrypt_Key"),
	}

	if GlobalOptions == nil {
		GlobalOptions = opts
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
