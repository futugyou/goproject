package verceltool

import (
	"context"
	"fmt"
	"html"
	"strings"

	"net/http"
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/services"

	"github.com/futugyou/extensions"
)

func AuthForVercel(w http.ResponseWriter, r *http.Request) bool {
	bearer := r.Header.Get("Authorization")
	_, err := Verify(r.Context(), bearer)
	if err != nil {
		w.Write([]byte("Auth error, you can get token again and try again."))
		w.WriteHeader(401)
		return false
	}

	return true
}

func Verify(ctx context.Context, authorization string) (*extensions.TokenVerifyResponse, error) {
	db_name := os.Getenv("db_name")
	authOptions := extensions.AuthOptions{
		AuthServerURL: os.Getenv("auth_server_url"),
		ClientID:      os.Getenv("client_id"),
		ClientSecret:  os.Getenv("client_secret"),
		Scopes:        os.Getenv("scopes"),
		RedirectURL:   os.Getenv("redirect_url"),
		AuthURL:       os.Getenv("auth_url"),
		TokenURL:      os.Getenv("token_url"),
		DbUrl:         os.Getenv("mongodb_url"),
		DbName:        &db_name,
	}

	oauthsvc := extensions.NewAuthMongoDBClient(authOptions)
	bearer := strings.ReplaceAll(authorization, "Bearer ", "")
	return oauthsvc.Verify(ctx, bearer)
}

func CheckAccountForVercel(w http.ResponseWriter, r *http.Request) bool {
	accountId := r.Header.Get("Account-Id")

	if len(accountId) == 0 || accountId == "aws-account-id-magic-code" {
		return true
	}

	accountService := services.NewAccountService()
	account := accountService.GetAccountByID(r.Context(), accountId)
	if account == nil || !account.Valid {
		w.Write([]byte("You may use a wrong 'Account-Id', please check."))
		w.WriteHeader(400)
		return false
	}

	err := awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	if err != nil {
		fmt.Fprintf(w, "The AWS SECRET associated with the account {%q} is incorrect.", html.EscapeString(accountId))
		w.WriteHeader(400)
		return false
	}

	return true
}
