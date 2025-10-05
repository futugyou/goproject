package api

import (
	"net/http"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/server"
	"github.com/futugyousuzu/identity/user"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store, err := session.Start(ctx, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")
		userstore := user.NewUserStore(ctx)
		user, err := userstore.Login(ctx, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		store.Set("LoggedInUserID", user.UserID)
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	server.OutputHTML(w, r, "login.html")
}
