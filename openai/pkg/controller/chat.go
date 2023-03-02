package controller

import (
	"openai/lib"
	"openai/pkg/services"

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
