package api

import (
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/devfeel/mapper"

	"github.com/futugyousuzu/go-openai-web/models"
	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/futugyou/extensions"
)

func Completions(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var completionModel models.CompletionModel
	err := json.NewDecoder(r.Body).Decode(&completionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&r)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	if !b {
		errString := ""
		for _, err := range valid.Errors {
			errString += (err.Key + " " + err.Message)
		}

		w.Write([]byte(errString))
		w.WriteHeader(500)
		return
	}

	completionService := services.CompletionService{}
	co := services.CompletionModel{}
	mapper.AutoMapper(&completionModel, &co)

	re := services.CreateCompletionRequest{
		CompletionModel: co,
	}

	result := completionService.CreateCompletion(re)
	body, _ := json.Marshal(result)
	w.Write(body)
	w.WriteHeader(200)
}
