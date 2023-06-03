package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var oauth_request_table = "oauth_request"

func (a *AuthService) saveToken(ctx context.Context, token *oauth2.Token) error {
	model := TokenModel{
		ID:           token.AccessToken,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	coll := a.Client.Database(a.Options.DbName).Collection(oauth_request_table)
	_, err := coll.InsertOne(ctx, model)
	return err
}

func (a *AuthService) getAuthRequestInfo(ctx context.Context, state string) AuthModel {
	var model AuthModel
	coll := a.Client.Database(a.Options.DbName).Collection(oauth_request_table)
	err := coll.FindOne(ctx, bson.D{{Key: "_id", Value: state}}).Decode(&model)
	if err != nil {
		panic(err)
	}

	return model
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func (a *AuthService) createAuthCodeURL(r *http.Request, config oauth2.Config) string {
	code_verifier := uuid.New().String()
	code_challenge := genCodeChallengeS256(code_verifier)
	code_challenge_method := "S256"
	state := strings.ReplaceAll(uuid.New().String(), "-", "")

	coll := a.Client.Database(a.Options.DbName).Collection(oauth_request_table)
	var model AuthModel = AuthModel{
		ID:                  state,
		CodeVerifier:        code_verifier,
		CodeChallenge:       code_challenge,
		CodeChallengeMethod: code_challenge_method,
		State:               state,
		RequestURI:          r.RequestURI,
		CreateAt:            time.Now(),
	}

	coll.InsertOne(r.Context(), model)

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

type AuthOptions struct {
	AuthServerURL string
	ClientID      string
	ClientSecret  string
	Scopes        string
	RedirectURL   string
	AuthURL       string
	TokenURL      string
	DbUrl         string
	DbName        string
}

type AuthService struct {
	Client      *mongo.Client
	Options     AuthOptions
	OauthConfig oauth2.Config
}

func NewAuthService(opts AuthOptions) *AuthService {
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(opts.DbUrl))
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
	return &AuthService{
		Client:      client,
		Options:     opts,
		OauthConfig: config,
	}
}

func (a *AuthService) RedirectToAuthorizationEndPoint(w http.ResponseWriter, r *http.Request) {
	authCodeURL := a.createAuthCodeURL(r, a.OauthConfig)
	http.Redirect(w, r, authCodeURL, http.StatusFound)
}

func (a *AuthService) Oauth2(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")
	state := r.Form.Get("state")

	if len(code) == 0 || len(state) == 0 {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	authModel := a.getAuthRequestInfo(r.Context(), state)
	token, err := a.OauthConfig.Exchange(r.Context(), code, oauth2.SetAuthURLParam("code_verifier", authModel.CodeVerifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.saveToken(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(token.AccessToken)
}

func (a *AuthService) VerifyToken(w http.ResponseWriter, r *http.Request, token *oauth2.Token) {
	set, err := jwk.Fetch(r.Context(), a.Options.AuthServerURL+".well-known/jwks.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tok, err := jwt.Parse([]byte(token.AccessToken), jwt.WithKeySet(set))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// "scope" can get from tok.PrivateClaims() or directly
	scope, _ := tok.Get("scope")
	fmt.Println(scope)

	fmt.Println(tok.Issuer())
	fmt.Println(tok.JwtID())
	fmt.Println(tok.Subject())
	for k, v := range tok.PrivateClaims() {
		fmt.Println(k, v)
	}

	//jws
	msg, _ := jws.Parse([]byte(token.AccessToken))
	for _, v := range msg.Signatures() {
		fmt.Println(v.ProtectedHeaders().KeyID())
		fmt.Println(v.ProtectedHeaders().Algorithm())
		fmt.Println(v.ProtectedHeaders().Get("x-example"))
	}
}
