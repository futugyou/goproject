package api

import (
	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
)

func Examples(w http.ResponseWriter, r *http.Request) {
	typestring := r.URL.Query().Get("type")
	exampleService := services.NewExampleService(createMongoDbCLient(), createRedisICLient())
	var result []services.ExampleModel
	if typestring == "custom" {
		result = exampleService.GetCustomExamples()
	} else {
		result = exampleService.GetSystemExamples()
	}

	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
