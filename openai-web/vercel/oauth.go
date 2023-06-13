package verceltool

import (
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"net/http"
	"os"

	oauthService "github.com/futugyousuzu/go-openai-web/oauth"
)

func AuthForVercel(w http.ResponseWriter, r *http.Request) bool {
	bearer := r.Header.Get("Authorization")
	authOptions := oauthService.AuthOptions{
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

	oauthsvc := oauthService.NewAuthService(authOptions)
	bearer = strings.ReplaceAll(bearer, "Bearer ", "")
	return oauthsvc.VerifyTokenString(w, r, bearer)
}
