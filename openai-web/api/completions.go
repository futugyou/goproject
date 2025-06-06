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
	b, err := valid.Valid(&completionModel)

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
		w.WriteHeader(400)
		return
	}

	completionService := services.NewCompletionService(createOpenAICLient())
	co := services.CompletionModel{}
	mapper.AutoMapper(&completionModel, &co)

	re := services.CreateCompletionRequest{
		CompletionModel: co,
	}

	result, err := completionService.CreateCompletion(r.Context(), re)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	w.Write(body)
	w.WriteHeader(200)
}
