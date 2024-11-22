package api

import (
	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
)

func Models(w http.ResponseWriter, r *http.Request) {
	modelService := services.NewModelService(createOpenAICLient(), createRedisICLient())
	result := modelService.GetAllModels(r.Context())
	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
