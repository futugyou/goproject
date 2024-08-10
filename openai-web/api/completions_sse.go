package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	"github.com/devfeel/mapper"

	"github.com/futugyousuzu/go-openai-web/models"
	"github.com/futugyousuzu/go-openai-web/services"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/futugyou/extensions"
)

func Completions_Sse(w http.ResponseWriter, r *http.Request) {
	if extensions.Cros(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	var buf []byte
	buf, _ = io.ReadAll(r.Body)

	var completionModel models.CompletionModel
	json.Unmarshal(buf, &completionModel)
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
		w.WriteHeader(500)
		return
	}

	completionService := services.CompletionService{}
	co := services.CompletionModel{}
	mapper.AutoMapper(&completionModel, &co)

	re := services.CreateCompletionRequest{
		CompletionModel: co,
	}

	result := completionService.CreateCompletionSSE(re)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for response := range result {
		message := strings.Join(response.Texts, ",")
		if len(message) == 0 {
			continue
		}

		message = url.QueryEscape(message)
		data := fmt.Sprintf("data: %s\n\n", message)
		w.Write([]byte(data))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	w.Write([]byte("data: [DONE]\n\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
