package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"
)

func ExamplesPost(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var buf []byte
	buf, _ = io.ReadAll(r.Body)
	var request services.ExampleModel
	json.Unmarshal(buf, &request)
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
