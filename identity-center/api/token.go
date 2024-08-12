package api

import (
	"net/http"

	"github.com/futugyousuzu/identity/server"

	"github.com/futugyou/extensions"
)

func Token(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	err := server.OAuthServer.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
