package api

import (
	"net/http"

	"github.com/futugyousuzu/identity/server"
	"github.com/futugyousuzu/identity/user"
	session "github.com/go-session/session/v3"
)

func Login(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
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
		userstore := user.NewUserStore()
		user, err := userstore.Login(r.Context(), username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		store.Set("LoggedInUserID", user.ID)
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	server.OutputHTML(w, r, "login.html")
}
