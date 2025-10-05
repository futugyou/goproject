package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"
	"os"

	"github.com/futugyousuzu/identity/token"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Query().Get("key") != "E86ZaWHjHnRqn" {
		w.WriteHeader(404)
		return
	}

	signed_key_id := os.Getenv("signed_key_id")
	token.NewJwksStore().CreateJwks(ctx, signed_key_id)
	w.WriteHeader(200)
}
