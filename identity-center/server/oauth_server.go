package server

import (
	"io"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"log"
	"os"

	"net/http"
	"net/http/httputil"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"

	"github.com/futugyousuzu/identity/mongo"
	"github.com/futugyousuzu/identity/user"
)

var OAuthServer *server.Server

func init() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	// manager.MustTokenStorage(store.NewMemoryTokenStore())
	mongodb_uri := os.Getenv("mongodb_url")
	mongodb_name := os.Getenv("db_name")

	manager.MapTokenStorage(
		mongo.NewTokenStore(mongo.NewConfig(
			mongodb_uri,
			mongodb_name,
		)),
	)
	// generate jwt access token
	signed_key_id := os.Getenv("signed_key_id")
	signed_key := os.Getenv("signed_key")
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate(signed_key_id, []byte(signed_key), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := mongo.NewClientStore(mongo.NewConfig(
		mongodb_uri,
		mongodb_name,
	))

	clientStore.Set(context.Background(), &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	OAuthServer = server.NewServer(server.NewConfig(), manager)
	OAuthServer.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		store := user.NewUserSore()
		user, err := store.Login(ctx, username, password)
		if err == nil {
			userID = user.Name
		}
		return
	})

	OAuthServer.SetUserAuthorizationHandler(UserAuthorizeHandler)

	OAuthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	OAuthServer.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
}

func OutputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func DumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
