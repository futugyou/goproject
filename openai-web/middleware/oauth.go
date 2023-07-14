package middleware

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	oauthService "github.com/futugyousuzu/identity/client"
)

type Options struct {
	AuthServerURL string
	ClientID      string
	ClientSecret  string
	Scopes        string
	RedirectURL   string
	AuthURL       string
	TokenURL      string
}

func OAuthConfig(opts *Options) web.FilterFunc {
	authOptions := oauthService.AuthOptions{
		AuthServerURL: opts.AuthServerURL,
		ClientID:      opts.ClientID,
		ClientSecret:  opts.ClientSecret,
		Scopes:        opts.Scopes,
		RedirectURL:   opts.RedirectURL,
		AuthURL:       opts.AuthURL,
		TokenURL:      opts.TokenURL,
		DbUrl:         os.Getenv("mongodb_url"),
		DbName:        os.Getenv("db_name"),
	}

	return func(ctx *context.Context) {
		oauthsvc := oauthService.NewAuthService(authOptions)

		if strings.HasPrefix(ctx.Request.RequestURI, fmt.Sprintf("/%s?code=", path.Base(opts.RedirectURL))) {
			oauthsvc.Oauth2(ctx.ResponseWriter, ctx.Request)
			return
		}

		// maybe swagger
		if !strings.HasPrefix(ctx.Request.RequestURI, "/api/") {
			return
		}

		auth := ctx.Request.Header.Get("Authorization")
		authorization := strings.ReplaceAll(auth, "Bearer ", "")
		if !strings.HasPrefix(auth, "Bearer ") || len(authorization) == 0 {
			ctx.ResponseWriter.WriteHeader(401)
			return
		}

		authorization = strings.ReplaceAll(authorization, "Bearer ", "")
		oauthsvc.VerifyTokenString(ctx.ResponseWriter, ctx.Request, authorization)
	}
}
