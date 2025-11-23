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
	VaultProjectUrl   string
	GetVaultsByIDs    string
	ShowVaultRaw      string
	CreateVault       string
	TriggerAuthKey    string
	ProjectWebhookUrl string
	ScreenshotAllow   string
	ScreenshotApiKey  string
	ScreenshotType    string
	GofileFolder      string
	GofileServer      string
	GofileToken       string
	VaultApiKey       string
}

func New() *Options {
	projectUrl := os.Getenv("VAULT_PROJECT_URL")
	opts := &Options{
		DBName:            os.Getenv("db_name"),
		MongoDBURL:        os.Getenv("mongodb_url"),
		QstashToken:       os.Getenv("QSTASH_TOKEN"),
		QstashDestination: os.Getenv("QSTASH_DESTINATION"),
		VaultProjectUrl:   projectUrl,
		GetVaultsByIDs:    GetEnvWithDefault(os.Getenv("VAULT_API_GET_BY_IDS"), createEndpoint(projectUrl, "/api/v1/vaults/by_ids")),
		ShowVaultRaw:      GetEnvWithDefault(os.Getenv("VAULT_API_SHOW_RAW"), createEndpoint(projectUrl, "/api/v1/vault/%s/show")),
		CreateVault:       GetEnvWithDefault(os.Getenv("VAULT_API_CREATE"), createEndpoint(projectUrl, "/api/v1/vault")),
		TriggerAuthKey:    os.Getenv("TRIGGER_AUTH_KEY"),
		ProjectWebhookUrl: os.Getenv("PROJECT_WEBHOOK_URL"),
		ScreenshotAllow:   GetEnvWithDefault(os.Getenv("SCREENSHOT_ALLOW"), "false"),
		ScreenshotApiKey:  os.Getenv("SCREENSHOT_API_KEY"),
		ScreenshotType:    GetEnvWithDefault(os.Getenv("SCREENSHOT_TYPE"), "Apiflash"),
		GofileFolder:      os.Getenv("GOFILE_FOLDER"),
		GofileServer:      os.Getenv("GOFILE_SERVER"),
		GofileToken:       os.Getenv("GOFILE_TOKEN"),
		VaultApiKey:       os.Getenv("VAULT_API_KEY"),
	}

	return opts
}

func createEndpoint(projectUrl, s string) string {
	if len(projectUrl) > 0 && projectUrl[len(projectUrl)-1] == '/' {
		projectUrl = projectUrl[:len(projectUrl)-1]
	}

	return projectUrl + s
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
