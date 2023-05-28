package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

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

var rawRequestUrl = ""

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
		if strings.HasPrefix(ctx.Request.RequestURI, fmt.Sprintf("/%s?code=", path.Base(opts.RedirectURL))) {
			ctx.Request.ParseForm()
			code := ctx.Request.Form.Get("code")
			state := ctx.Request.Form.Get("state")

			if len(code) == 0 || len(state) == 0 {
				http.Error(ctx.ResponseWriter, "State invalid", http.StatusBadRequest)
				return
			}

			token, err := config.Exchange(ctx.Request.Context(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
			if err != nil {
				http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx.ResponseWriter.Header().Set("Authorization", token.TokenType+" "+token.AccessToken)
			http.Redirect(ctx.ResponseWriter, ctx.Request, rawRequestUrl, http.StatusFound)
			return
		}

		if !strings.HasPrefix(ctx.Request.RequestURI, "/api/") {
			return
		}

		authorization := ctx.Request.Header.Get("Authorization")
		fmt.Println(ctx.Request.RequestURI, authorization)
		if len(authorization) == 0 {
			rawRequestUrl = ctx.Request.RequestURI
			u := config.AuthCodeURL("xyz",
				oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
				oauth2.SetAuthURLParam("code_challenge_method", "S256"))
			http.Redirect(ctx.ResponseWriter, ctx.Request, u, http.StatusFound)

			return
		}
	}
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
