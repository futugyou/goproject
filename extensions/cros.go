package extensions

import "net/http"

func Cors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "traceparent, Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Token, Content-Length, Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func CorsWithOrigin(w http.ResponseWriter, r *http.Request, origin string) bool {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "traceparent, Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Token, Content-Length, Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func CorsWithMultipleOrigin(w http.ResponseWriter, r *http.Request, origins []string) bool {
	allowedOrigins := make(map[string]bool)
	for _, org := range origins {
		allowedOrigins[org] = true
	}

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "traceparent, Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Token, Content-Length, Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func CorsAll(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}
