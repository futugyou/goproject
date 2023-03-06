package controllers

import (
	"openai/lib"
	"openai/services"

	"github.com/beego/beego/v2/server/web"
)

type ChatController struct {
	web.Controller
}

func (c *ChatController) ListModel(request lib.CreateChatCompletionRequest) {
	chatService := services.ChatService{}
	result := chatService.CreateChatCompletion(request)
	c.Ctx.JSONResp(result)
}
