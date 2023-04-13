package controllers

import (
	"encoding/json"
	"fmt"  
	"net/url"
	"strings"

	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about completion
type CompletionController struct {
	web.Controller
}

// @Title Create Completion With SSE
// @Description create chat
// @Param	body		body 	lib.CreateChatCompletionRequest	true		"body for create completion content"
// @Success 200 {object} 	lib.CreateChatCompletionResponse
// @router / [post]
func (c *CompletionController) CreateCompletionWithSSE() {
	completionService := services.CompletionService{}
	var completion services.CreateCompletionRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &completion)
	result := completionService.CreateCompletionSSE(completion)

	var rw = c.Ctx.ResponseWriter
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	for response := range result {
		message := strings.Join(response.Texts, ",")
		if len(message) == 0 {
			continue
		}

		message  = url.QueryEscape(message)
		data := fmt.Sprintf("data: %s\n\n", message)
		rw.Write([]byte(data))
		rw.Flush()
	}
	rw.Write([]byte("data: [DONE]\n\n"))
	rw.Flush()
}
