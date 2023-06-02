package server

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/lestrrat-go/jwx/v2/jwa"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/mongo"
	sessionstore "github.com/futugyousuzu/identity/session"
	assets "github.com/futugyousuzu/identity/static"
	"github.com/futugyousuzu/identity/token"
)

var OAuthServer *server.Server

func init() {
	mongodb_uri := os.Getenv("mongodb_url")
	mongodb_name := os.Getenv("db_name")
	signed_key_id := os.Getenv("signed_key_id")
	signed_key := os.Getenv("signed_key")

	token.NewJwksStore().CreateJwks(context.Background(), signed_key_id)

	session.InitManager(
		session.SetStore(sessionstore.NewStore(mongodb_uri, mongodb_name, "session")),
	)

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	// manager.MustTokenStorage(store.NewMemoryTokenStore())

	manager.MapTokenStorage(
		mongo.NewTokenStore(mongo.NewConfig(
			mongodb_uri,
			mongodb_name,
		)),
	)
	// generate jwt access token
	manager.MapAccessGenerate(token.NewJWTAccessGenerate(signed_key_id, []byte(signed_key), jwa.RS256))

	clientStore := mongo.NewClientStore(mongo.NewConfig(
		mongodb_uri,
		mongodb_name,
	))

	initClient(clientStore)
	manager.MapClientStorage(clientStore)

	OAuthServer = server.NewServer(server.NewConfig(), manager)
	OAuthServer.SetPasswordAuthorizationHandler(PasswordAuthorizationHandler)

	OAuthServer.SetUserAuthorizationHandler(UserAuthorizeHandler)
	OAuthServer.SetAuthorizeScopeHandler(AuthorizeScopeHandler)

	OAuthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	OAuthServer.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
}

func OutputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	a := &assets.Assets
	file, err := a.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, filename, fi.ModTime(), file.(io.ReadSeeker))
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
