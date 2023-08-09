package verceltool

import (
	"fmt"
	"strings"

	"net/http"
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/services"
	oauthService "github.com/futugyousuzu/identity/client"
)

func AuthForVercel(w http.ResponseWriter, r *http.Request) bool {
	accountId := r.Header.Get("Account-Id")
	bearer := r.Header.Get("Authorization")

	if len(accountId) == 0 || len(bearer) == 0 {
		w.Write([]byte("This request NEED both 'Authorization' and 'Account-Id' in header."))
		w.WriteHeader(400)
		return false
	}

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
	verifyResult := oauthsvc.VerifyTokenString(w, r, bearer)
	if !verifyResult {
		w.Write([]byte("Auth error, you can get token again and try again."))
		w.WriteHeader(401)
		return false
	}

	if accountId == "aws-account-id-magic-code" {
		return true
	}

	accountService := services.NewAccountService()
	account := accountService.GetAccountByID(accountId)
	if account == nil {
		w.Write([]byte("You may use a wrong 'Account-Id', please check."))
		w.WriteHeader(400)
		return false
	}

	err := awsenv.CfgForVercel(account.AccessKeyId, account.SecretAccessKey)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("The AWS SECRET associated with the account {%s} is incorrect.", accountId)))
		w.WriteHeader(400)
		return false
	}

	return true
}
