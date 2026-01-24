package main

import (
	"fmt"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/api"
)

func main() {
	http.HandleFunc("/light", api.AguiHandler)

	fmt.Println("server started, address: http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("server error: %s\n", err)
	}
}
