package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"
)

func Edits(w http.ResponseWriter, r *http.Request) {
	if verceltool.CrosForVercel(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var editsRequest services.CreateEditsRequest
	var buf []byte
	buf, _ = io.ReadAll(r.Body)
	json.Unmarshal(buf, &editsRequest)

	completionService := services.EditService{}
	result := completionService.CreateEdit(editsRequest)
	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
