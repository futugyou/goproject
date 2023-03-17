package controllers

import (
	"encoding/json"

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
	chatService := services.ChatService{}
	var chat lib.CreateChatCompletionRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &chat)
	result := chatService.CreateChatCompletion(chat)
	c.Ctx.JSONResp(result)
}
