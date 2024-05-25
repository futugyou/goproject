package internal

import (
	"net/http"
)

func CorsForVercel(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Token, Content-Length, Access-Control-Allow-Headers, Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}
