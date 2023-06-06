package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"
	"os"

	"github.com/futugyousuzu/go-openai-web/oauth"
)

func auth(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	authOptions := oauth.AuthOptions{
		AuthServerURL: os.Getenv("auth_server_url"),
		ClientID:      os.Getenv("client_id"),
		ClientSecret:  os.Getenv("client_secret"),
		Scopes:        os.Getenv("scopes"),
		RedirectURL:   os.Getenv("redirect_url"),
		AuthURL:       os.Getenv("auth_url"),
		TokenURL:      os.Getenv("token_url"),
		DbUrl:         os.Getenv("mongodb_url"),
		DbName:        os.Getenv("db_name"),
	}

	oauthsvc := oauth.NewAuthService(authOptions)
	oauthsvc.VerifyTokenString(w, r, bearer)
}
