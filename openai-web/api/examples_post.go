package api

import (
	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/futugyou/extensions"
)

func ExamplesPost(w http.ResponseWriter, r *http.Request) {
	if extensions.Cros(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var request services.ExampleModel

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exampleService := services.ExampleService{}
	if len(request.Key) == 0 {
		w.Write([]byte("errors"))
		w.WriteHeader(500)
		return
	}

	typestring := r.URL.Query().Get("type")
	if typestring == "custom" {
		exampleService.CreateCustomExample(request)
	} else {
		exampleService.CreateSystemExample(request)
	}

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
