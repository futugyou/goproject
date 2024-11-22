package api

import (
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/futugyou/extensions"
)

func ExamplesReset(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	exampleService := services.NewExampleService(createMongoDbCLient(r.Context()), createRedisICLient())
	exampleService.Reset(r.Context())

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
