package controllers

import (
	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Completion
type CompletionController struct {
	web.Controller
}

// @Title get settings
// @Param settingName path string true "setting name"
// @Success 200 {object} 	services.CompletionModel
// @router /:settingName [get]
func (c *CompletionController) Setting() {
	settingName := c.Ctx.Input.Param(":settingName")
	chatService := services.CompletionService{}
	result := chatService.GetExampleSettings(settingName)
	c.Ctx.JSONResp(result)
}
