package oauth

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
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
	// scopes := make([]string, 0)
	// json.Unmarshal([]byte(opts.Scopes), &scopes)

	// config := oauth2.Config{
	// 	ClientID:     opts.ClientID,
	// 	ClientSecret: opts.ClientSecret,
	// 	Scopes:       scopes,
	// 	RedirectURL:  opts.RedirectURL,
	// 	Endpoint: oauth2.Endpoint{
	// 		AuthURL:  opts.AuthServerURL + opts.AuthURL,
	// 		TokenURL: opts.AuthServerURL + opts.TokenURL,
	// 	},
	// }

	return func(ctx *context.Context) {
		if !strings.HasPrefix(ctx.Request.RequestURI, "/api/") {
			return
		}

		authorization := ctx.Request.Header.Get("Authorization")

		if len(authorization) == 0 {
			ctx.ResponseWriter.WriteHeader(401)
			return
		}
	}
}
