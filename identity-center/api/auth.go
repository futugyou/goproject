package api

import (
	_ "github.com/joho/godotenv/autoload"

	"net/http"

	session "github.com/go-session/session/v3"

	"github.com/futugyousuzu/identity/server"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	store, err := session.Start(ctx, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	base := server.GetBaseUrl(r)
	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", base+"/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	// server.OutputHTML(w, r, "auth.html")
	
	data := map[string]interface{}{ 
		"Base":    base,
	}

	server.OutputHTMLWithData(w, r, "auth.html", data)
}
