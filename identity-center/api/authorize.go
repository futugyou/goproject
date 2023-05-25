package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/futugyousuzu/identity/server"
	session "github.com/go-session/session/v3"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values = make(url.Values)
	if dic, ok := store.Get("ReturnUri"); ok {
		for k, v := range dic.(map[string]interface{}) {
			vv := v.([]interface{})
			s := make([]string, len(vv))
			for i, v := range vv {
				s[i] = fmt.Sprint(v)
			}
			form[k] = s
		}
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = server.OAuthServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
