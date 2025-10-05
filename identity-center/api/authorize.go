package api

import (
	"fmt"
	"net/http"
	"net/url"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/server"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store, err := session.Start(ctx, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if dic, ok := store.Get("ReturnUri"); ok {
		form = make(url.Values)
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
