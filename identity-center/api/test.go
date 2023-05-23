package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/futugyousuzu/identity/server"
)

func Test(w http.ResponseWriter, r *http.Request) {
	token, err := server.OAuthServer.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}
