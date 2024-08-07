package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/server"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	server.OutputHTML(w, r, "auth.html")
}
