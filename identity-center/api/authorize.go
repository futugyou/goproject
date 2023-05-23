package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/futugyousuzu/identity/server"
	"github.com/go-session/session"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		fmt.Print("sss")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = server.OAuthServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		fmt.Print("bbbb")

		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
