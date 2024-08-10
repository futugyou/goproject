package api

import (
	"encoding/json"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/futugyou/extensions"
)

func Edits(w http.ResponseWriter, r *http.Request) {
	if extensions.Cros(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var editsRequest services.CreateEditsRequest
	err := json.NewDecoder(r.Body).Decode(&editsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	completionService := services.EditService{}
	result := completionService.CreateEdit(editsRequest)
	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
