package options

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Options struct {
	MongoDBURL              string
	DBName                  string
	QueryDBName             string
	QueryMongoDBURL         string
	CircleCIOrgSlug         string
	HCPClientID             string
	HCPClientSecret         string
	HCPOrganizationID       string
	HCPProjectID            string
	HCPAppName              string
	TFCOrganizationToken    string
	TFCToken                string
	TFCAPIBaseURL           string
	TFCOrg                  string
	TFCWorkspace            string
	EncryptKey              string
	GithubToken             string
	VercelToken             string
	CircleCIToken           string
	QstashToken             string
	QstashDestination       string
	QstashCurrentSigningKey string
	QstashNextSigningKey    string
	AuthServerURL           string
	ClientID                string
	ClientSecret            string
	Scopes                  string
	RedirectURL             string
	AuthURL                 string
	TokenURL                string
	RedisURL                string
	ProjectWebhookURL       string
	ProjectURL              string
	TriggerAuthKey          string
	ScreenshotAPIKey        string
	ScreenshotType          string
	ScreenshotAllow         string
	GoFileToken             string
	GoFileServer            string
	GoFileFolder            string
	EventPublisher          string
	StoreType               string
}

func New(origArgs []string) (*Options, error) {
	opts := &Options{
		MongoDBURL:              os.Getenv("mongodb_url"),
		DBName:                  os.Getenv("db_name"),
		QueryDBName:             os.Getenv("query_db_name"),
		QueryMongoDBURL:         os.Getenv("query_mongodb_url"),
		CircleCIOrgSlug:         os.Getenv("CIRCLECI_ORG_SLUG"),
		HCPClientID:             os.Getenv("HCP_CLIENT_ID"),
		HCPClientSecret:         os.Getenv("HCP_CLIENT_SECRET"),
		HCPOrganizationID:       os.Getenv("HCP_ORGANIZATION_ID"),
		HCPProjectID:            os.Getenv("HCP_PROJECT_ID"),
		HCPAppName:              os.Getenv("HCP_APP_NAME"),
		TFCOrganizationToken:    os.Getenv("TFC_ORGANIZATION_TOKEN"),
		TFCToken:                os.Getenv("TFC_TOKEN"),
		TFCAPIBaseURL:           os.Getenv("TFC_APIBASEURL"),
		TFCOrg:                  os.Getenv("TFC_ORG"),
		TFCWorkspace:            os.Getenv("TFC_WORKSPACE"),
		EncryptKey:              os.Getenv("Encrypt_Key"),
		GithubToken:             os.Getenv("GITHUB_TOKEN"),
		VercelToken:             os.Getenv("VERCEL_TOKEN"),
		CircleCIToken:           os.Getenv("CIRCLECI_TOKEN"),
		QstashToken:             os.Getenv("QSTASH_TOKEN"),
		QstashDestination:       os.Getenv("QSTASH_DESTINATION"),
		QstashCurrentSigningKey: os.Getenv("QSTASH_CURRENT_SIGNING_KEY"),
		QstashNextSigningKey:    os.Getenv("QSTASH_NEXT_SIGNING_KEY"),
		AuthServerURL:           os.Getenv("auth_server_url"),
		ClientID:                os.Getenv("client_id"),
		ClientSecret:            os.Getenv("client_secret"),
		Scopes:                  os.Getenv("scopes"),
		RedirectURL:             os.Getenv("redirect_url"),
		AuthURL:                 os.Getenv("auth_url"),
		TokenURL:                os.Getenv("token_url"),
		RedisURL:                os.Getenv("REDIS_URL"),
		ProjectWebhookURL:       os.Getenv("PROJECT_WEBHOOK_URL"),
		ProjectURL:              os.Getenv("PROJECT_URL"),
		TriggerAuthKey:          os.Getenv("TRIGGER_AUTH_KEY"),
		ScreenshotAPIKey:        os.Getenv("SCREENSHOT_API_KEY"),
		ScreenshotType:          os.Getenv("SCREENSHOT_TYPE"),
		ScreenshotAllow:         os.Getenv("SCREENSHOT_ALLOW"),
		GoFileToken:             os.Getenv("GOFILE_TOKEN"),
		GoFileServer:            os.Getenv("GOFILE_SERVER"),
		GoFileFolder:            os.Getenv("GOFILE_FOLDER"),
		EventPublisher:          GetEnvWithDefault("EVENT_PUBLISHER", "qstash"),
		StoreType:               GetEnvWithDefault("STORE_TYPE", "mongo"),
	}

	return opts, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
