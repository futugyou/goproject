package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"golang.org/x/oauth2"
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
	scopes := make([]string, 0)
	json.Unmarshal([]byte(opts.Scopes), &scopes)

	config := oauth2.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		Scopes:       scopes,
		RedirectURL:  opts.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  opts.AuthServerURL + opts.AuthURL,
			TokenURL: opts.AuthServerURL + opts.TokenURL,
		},
	}

	return func(ctx *context.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		if len(authorization) == 0 {
			u := config.AuthCodeURL("xyz",
				oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
				oauth2.SetAuthURLParam("code_challenge_method", "S256"))
			http.Redirect(ctx.ResponseWriter, ctx.Request, u, http.StatusFound)
		}
	}
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
