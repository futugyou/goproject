package client

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
)

var oauth_request_table = "oauth_request"

func (a *AuthService) saveToken(ctx context.Context, token *oauth2.Token) error {
	model := &TokenModel{
		ID:           token.AccessToken,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	return a.OAuthTokenRepository.Insert(ctx, model)
}

func (a *AuthService) getAuthRequestInfo(ctx context.Context, state string) *AuthModel {
	model, err := a.OAuthRequestRepository.Get(ctx, state)
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

	var model AuthModel = AuthModel{
		ID:                  state,
		CodeVerifier:        code_verifier,
		CodeChallenge:       code_challenge,
		CodeChallengeMethod: code_challenge_method,
		State:               state,
		RequestURI:          r.RequestURI,
		CreateAt:            time.Now(),
	}

	a.OAuthRequestRepository.Insert(r.Context(), &model)

	return config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", code_challenge),
		oauth2.SetAuthURLParam("code_challenge_method", code_challenge_method),
		oauth2.AccessTypeOffline)
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
	Options                AuthOptions
	OauthConfig            oauth2.Config
	OAuthRequestRepository IOAuthRequestRepository
	OAuthTokenRepository   IOAuthTokenRepository
}

func NewAuthService(opts AuthOptions) *AuthService {
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

func (a *AuthService) VerifyToken(w http.ResponseWriter, r *http.Request, token *oauth2.Token) bool {
	return a.VerifyTokenString(w, r, token.AccessToken)
}

func (a *AuthService) VerifyTokenString(w http.ResponseWriter, r *http.Request, authorization string) bool {
	set, err := jwk.Fetch(r.Context(), a.Options.AuthServerURL+".well-known/jwks.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	tok, err := jwt.Parse([]byte(authorization), jwt.WithKeySet(set))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
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
	msg, _ := jws.Parse([]byte(authorization))
	for _, v := range msg.Signatures() {
		fmt.Println(v.ProtectedHeaders().KeyID())
		fmt.Println(v.ProtectedHeaders().Algorithm())
		fmt.Println(v.ProtectedHeaders().Get("x-example"))
	}

	return true
}
