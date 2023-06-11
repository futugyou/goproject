package main

import (
	_ "github.com/joho/godotenv/autoload"

	_ "github.com/futugyousuzu/identity/server"

	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/futugyousuzu/identity/api"
	"github.com/futugyousuzu/identity/apiraw"
	"github.com/futugyousuzu/identity/middleware"
	"github.com/futugyousuzu/identity/scheduler"
)

var (
	portvar int
)

func init() {
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()
	scheduler.JwksGenerate()
	http.HandleFunc("/login", middleware.Cors(api.Login))
	http.HandleFunc("/auth", middleware.Cors(api.Auth))
	http.HandleFunc("/authorize", middleware.Cors(api.Authorize))
	http.HandleFunc("/token", middleware.Cors(apiraw.Token))
	http.HandleFunc("/test", middleware.Cors(api.Test))
	http.HandleFunc("/.well-known/jwks.json", middleware.Cors(api.Jwks))
	http.HandleFunc("/", middleware.Cors(api.Jwks))

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/token")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}
