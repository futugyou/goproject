package middleware

import "net/http"

func Cors(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Origin,Authorization,Access-Control-Allow-Headers,Content-Type，Token")
		w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers,Token,Content-Length,Access-Control-Allow-Headers,Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next(w, r)
	})
}
