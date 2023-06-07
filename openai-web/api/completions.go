package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
	"github.com/devfeel/mapper"

	"github.com/futugyousuzu/go-openai-web/models"
	// "github.com/futugyousuzu/go-openai-web/oauth"
	"github.com/futugyousuzu/go-openai-web/services"
)

func Completions(w http.ResponseWriter, r *http.Request) {
	// oauth.AuthForVercel(w, r)

	var buf []byte
	buf, _ = io.ReadAll(r.Body)

	var completionModel models.CompletionModel
	json.Unmarshal(buf, &completionModel)
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
