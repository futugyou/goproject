package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/oauth2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

var oauth_request_table = "oauth_request"

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
		fmt.Println(ctx.Request.RequestURI, fmt.Sprintf("/%s?code=", path.Base(opts.RedirectURL)))
		if strings.HasPrefix(ctx.Request.RequestURI, fmt.Sprintf("/%s?code=", path.Base(opts.RedirectURL))) {
			ctx.Request.ParseForm()
			code := ctx.Request.Form.Get("code")
			state := ctx.Request.Form.Get("state")

			if len(code) == 0 || len(state) == 0 {
				http.Error(ctx.ResponseWriter, "State invalid", http.StatusBadRequest)
				return
			}

			authModel := getAuthRequestInfo(ctx, state)
			token, err := config.Exchange(ctx.Request.Context(), code, oauth2.SetAuthURLParam("code_verifier", authModel.CodeVerifier))
			if err != nil {
				http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			err = saveToken(ctx, token)
			if err != nil {
				http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			// verification
			set, err := jwk.Fetch(ctx.Request.Context(), opts.AuthServerURL+".well-known/jwks.json")
			if err != nil {
				http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			tok, err := jwt.Parse([]byte(token.AccessToken), jwt.WithKeySet(set))
			if err != nil {
				http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			scope, _ := tok.Get("scope")
			fmt.Println(tok.Issuer())
			fmt.Println(scope)
			fmt.Println(tok.JwtID())
			fmt.Println(tok.Subject())
			for k, v := range tok.PrivateClaims() {
				fmt.Println(k, v)
			}

			return
		}

		if !strings.HasPrefix(ctx.Request.RequestURI, "/api/") {
			return
		}

		authorization := ctx.Request.Header.Get("Authorization")

		if len(authorization) == 0 {
			authCodeURL := createAuthCodeURL(ctx, config)
			http.Redirect(ctx.ResponseWriter, ctx.Request, authCodeURL, http.StatusFound)

			return
		}
	}
}

func saveToken(ctx *context.Context, token *oauth2.Token) error {
	model := TokenModel{
		ID:           token.AccessToken,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(ctx.Request.Context(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx.Request.Context()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database(db_name).Collection(oauth_request_table)
	_, err = coll.InsertOne(ctx.Request.Context(), model)
	return err
}

func getAuthRequestInfo(ctx *context.Context, state string) AuthModel {
	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(ctx.Request.Context(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx.Request.Context()); err != nil {
			panic(err)
		}
	}()

	var model AuthModel
	coll := client.Database(db_name).Collection(oauth_request_table)
	err = coll.FindOne(ctx.Request.Context(), bson.D{{Key: "_id", Value: state}}).Decode(&model)
	if err != nil {
		panic(err)
	}

	return model
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func createAuthCodeURL(ctx *context.Context, config oauth2.Config) string {
	uri := os.Getenv("mongodb_url")
	db_name := os.Getenv("db_name")
	client, err := mongo.Connect(ctx.Request.Context(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx.Request.Context()); err != nil {
			panic(err)
		}
	}()

	code_verifier := uuid.New().String()
	code_challenge := genCodeChallengeS256(code_verifier)
	code_challenge_method := "S256"
	state := strings.ReplaceAll(uuid.New().String(), "-", "")

	coll := client.Database(db_name).Collection(oauth_request_table)
	var model AuthModel = AuthModel{
		ID:                  state,
		CodeVerifier:        code_verifier,
		CodeChallenge:       code_challenge,
		CodeChallengeMethod: code_challenge_method,
		State:               state,
		RequestURI:          ctx.Request.RequestURI,
		CreateAt:            time.Now(),
	}

	coll.InsertOne(ctx.Request.Context(), model)

	return config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", code_challenge_method),
		oauth2.AccessTypeOffline)
}

type TokenModel struct {
	ID           string    `bson:"_id"`
	AccessToken  string    `bson:"access_token"`
	TokenType    string    `bson:"token_type"`
	RefreshToken string    `bson:"refresh_token"`
	Expiry       time.Time `bson:"expiry"`
}

type AuthModel struct {
	ID                  string    `bson:"_id"`
	CodeVerifier        string    `bson:"code_verifier"`
	CodeChallenge       string    `bson:"code_challenge"`
	CodeChallengeMethod string    `bson:"code_challenge_method"`
	State               string    `bson:"state"`
	RequestURI          string    `bson:"request_uri"`
	CreateAt            time.Time `bson:"create_at"`
}
