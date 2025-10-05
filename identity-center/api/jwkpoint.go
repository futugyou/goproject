package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	"github.com/futugyousuzu/identity/token"
)

func Jwks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store := token.NewJwksStore()
	result, _ := store.GetPublicJwksList(ctx)
	w.Write([]byte(result))
	w.WriteHeader(200)
}
