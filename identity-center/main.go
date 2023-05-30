package main

import (
	_ "github.com/joho/godotenv/autoload"

	_ "github.com/futugyousuzu/identity/server"

	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/futugyousuzu/identity/api"
)

var (
	portvar int
)

func init() {
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()

	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/auth", api.Auth)
	http.HandleFunc("/authorize", api.Authorize)
	http.HandleFunc("/token", api.Token)
	http.HandleFunc("/test", api.Test)
	http.HandleFunc("/jwks", api.Jwks)

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/token")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}
