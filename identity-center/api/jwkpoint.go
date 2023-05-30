package api

import (
	"github.com/futugyousuzu/identity/server"
	_ "github.com/joho/godotenv/autoload"

	"net/http"
)

func Jwks(w http.ResponseWriter, r *http.Request) {
	store := server.NewJwksStore()
	result, _ := store.GetJwksList(r.Context())
	w.Write([]byte(result))
	w.WriteHeader(200)
}
