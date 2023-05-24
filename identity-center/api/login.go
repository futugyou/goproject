package api

import (
	"fmt"
	"net/http"

	"github.com/futugyousuzu/identity/server"
	"github.com/futugyousuzu/identity/user"
	"github.com/go-session/session"
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
		userstore := user.NewUserSore()
		fmt.Println(username, password)
		user, err := userstore.Login(r.Context(), username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		store.Set("LoggedInUserID", user.Name)
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	server.OutputHTML(w, r, "static/login.html")
}
