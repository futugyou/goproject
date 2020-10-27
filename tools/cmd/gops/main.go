package main

import (
	"log"
	"net/http"

	"github.com/google/gops/agent"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal("agent listen err : %v", err)
	}
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("golang projecct"))
	})
	_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}
