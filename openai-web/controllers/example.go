package controllers

import (
	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Examples
type ExampleController struct {
	web.Controller
}

// @Title get examples
// @Success 200 {object} 	services.ExampleModel
// @router / [get]
func (c *ExampleController) ExampleDetail() {
	chatService := services.ExampleService{}
	result := chatService.GetExampleSettings()
	c.Ctx.JSONResp(result)
}
