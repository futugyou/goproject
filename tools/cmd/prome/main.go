package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("go project"))
	})
	_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}
