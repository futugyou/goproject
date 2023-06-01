package api

import (
	"github.com/futugyousuzu/identity/token"
	_ "github.com/joho/godotenv/autoload"

	"net/http"
)

func Jwks(w http.ResponseWriter, r *http.Request) {
	store := token.NewJwksStore()
	result, _ := store.GetPublicJwksList(r.Context())
	w.Write([]byte(result))
	w.WriteHeader(200)
}
