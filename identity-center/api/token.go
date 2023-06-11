package api

import (
	"net/http"

	"github.com/futugyousuzu/identity/apiraw"
	"github.com/futugyousuzu/identity/middleware"
)

func Token(w http.ResponseWriter, r *http.Request) {
	middleware.Cors(apiraw.Token)
}
