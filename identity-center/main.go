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
	dumpvar   bool
	idvar     string
	secretvar string
	domainvar string
	portvar   int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.StringVar(&idvar, "i", "222222", "The client id being passed in")
	flag.StringVar(&secretvar, "s", "22222222", "The client secret being passed in")
	flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()

	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/auth", api.Auth)
	http.HandleFunc("/authorize", api.Authorize)
	http.HandleFunc("/token", api.Token)
	http.HandleFunc("/test", api.Test)

	log.Printf("Server is running at %d port.\n", portvar)
	log.Printf("Point your OAuth client Auth endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/authorize")
	log.Printf("Point your OAuth client Token endpoint to %s:%d%s", "http://localhost", portvar, "/oauth/token")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portvar), nil))
}
