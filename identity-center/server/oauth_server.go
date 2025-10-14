package server

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"text/template"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

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

	mongostore "github.com/futugyousuzu/identity/mongo"
	sessionstore "github.com/futugyousuzu/identity/session"
	assets "github.com/futugyousuzu/identity/static"
	"github.com/futugyousuzu/identity/token"
)

var OAuthServer *server.Server

func maskMongoURI(uri string) string {
	reUser := regexp.MustCompile(`(?m)(mongodb\+srv://)([^@]+)@`)
	uri = reUser.ReplaceAllString(uri, "${1}****@")
	reCluster := regexp.MustCompile(`(cluster\d*\.)[^.]+`)
	uri = reCluster.ReplaceAllString(uri, "${1}****")
	return uri
}

func init() {
	mongodb_uri := os.Getenv("mongodb_url")
	mongodb_name := os.Getenv("db_name")
	signed_key_id := os.Getenv("signed_key_id")
	signed_key := os.Getenv("signed_key")

	mask_url := maskMongoURI(mongodb_uri)
	fmt.Printf("masked mongo url is %s, mongodb name is %s .\n", mask_url, mongodb_name)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_uri))
	if err != nil {
		fmt.Println("mongo connect failed:", err)
		return
	}

	// ping
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("mongo ping failed:", err)
		return
	}

	token.NewJwksStoreWithMongoClient(client).CreateJwks(ctx, signed_key_id)

	session.InitManager(
		session.SetStore(sessionstore.NewStoreWithMongoClient(client, mongodb_name, "session")),
	)

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	// manager.MustTokenStorage(store.NewMemoryTokenStore())

	manager.MapTokenStorage(
		mongostore.NewTokenStoreWithclient(client, mongodb_name),
	)
	// generate jwt access token
	manager.MapAccessGenerate(token.NewJWTAccessGenerate(signed_key_id, []byte(signed_key), jwa.RS256))

	clientStore := mongostore.NewClientStoreWithclient(client, mongodb_name)

	initClient(clientStore)
	manager.MapClientStorage(clientStore)

	authServerConfig := server.NewConfig()
	authServerConfig.ForcePKCE = true
	OAuthServer = server.NewServer(authServerConfig, manager)
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

func OutputHTMLWithData(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
	a := &assets.Assets
	file, err := a.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl, err := template.New(filename).Parse(string(content))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, data)
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

func GetBaseUrl(r *http.Request) string {
	forwarded_hosts_string := os.Getenv("forwarded_hosts")
	forwarded_hosts := make([]string, 0)
	json.Unmarshal([]byte(forwarded_hosts_string), &forwarded_hosts)
	if slices.Contains(forwarded_hosts, r.Header.Get("X-Forwarded-Host")) {
		return "/identity"
	}

	return ""
}
