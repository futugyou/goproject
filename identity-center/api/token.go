package api

import (
	"net/http"

	"github.com/futugyousuzu/identity/server"
)

func Token(w http.ResponseWriter, r *http.Request) {
	err := server.OAuthServer.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
