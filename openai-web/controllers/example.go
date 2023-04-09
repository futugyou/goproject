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
// @Param exampleName path string true "example name"
// @Success 200 {object} 	services.CompletionModel
// @router /:exampleName [get]
func (c *ExampleController) ExampleDetail() {
	exampleName := c.Ctx.Input.Param(":exampleName")
	chatService := services.ExampleService{}
	result := chatService.GetExampleSettings(exampleName)
	c.Ctx.JSONResp(result)
}
