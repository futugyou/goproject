package api

import (
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"
)

func ExamplesReset(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	exampleService := services.ExampleService{}
	exampleService.Reset()

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
