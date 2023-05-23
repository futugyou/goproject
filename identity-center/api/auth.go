package api

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyousuzu/identity/server"

	"net/http"

	"github.com/go-session/session"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	server.OutputHTML(w, r, "static/auth.html")
}
