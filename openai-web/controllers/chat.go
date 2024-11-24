package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"

	lib "github.com/futugyousuzu/go-openai"

	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Chat
type ChatController struct {
	web.Controller
}

// @Title CreateChat
// @Description create chat
// @Param	body		body 	lib.CreateChatCompletionRequest	true		"body for create chat content"
// @Success 200 {object} 	lib.CreateChatCompletionResponse
// @router / [post]
func (c *ChatController) CreateChat(request lib.CreateChatCompletionRequest) {
	chatService := services.NewChatService(createOpenAICLient())
	var chat lib.CreateChatCompletionRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &chat)
	result := chatService.CreateChatCompletion(c.Ctx.Request.Context(), chat)
	c.Ctx.JSONResp(result)
}

// @Title Create Chat With SSE
// @Description Create Chat Stream
// @Param	body		body 	services.CreateChatRequest	true		"body for create Chat content"
// @router /sse [post]
func (c *ChatController) CreateChatWithSSE() {
	chatService := services.NewChatService(createOpenAICLient())
	var request services.CreateChatRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)

	result := chatService.CreateChatSSE(c.Ctx.Request.Context(), request)

	var rw = c.Ctx.ResponseWriter
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	for response := range result {
		if len(response.Messages) == 0 {
			continue
		}

		message := ""
		for i := 0; i < len(response.Messages); i++ {
			content := response.Messages[i].Content
			if len(content) > 0 {
				message += content
			}
		}

		message = url.QueryEscape(message)
		data := fmt.Sprintf("data: %s\n\n", message)
		rw.Write([]byte(data))
		rw.Flush()
	}

	rw.Write([]byte("data: [DONE]\n\n"))
	rw.Flush()
}
