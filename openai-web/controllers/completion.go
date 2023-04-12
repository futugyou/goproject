package controllers

import (
	"encoding/json"

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

	for response := range result {
		var rw = c.Ctx.ResponseWriter
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
		rw.Header().Set("Cache-Control", "no-cache")
		data, _ := json.Marshal(&response)
		rw.Write(data)
		rw.Flush()
	}
}
