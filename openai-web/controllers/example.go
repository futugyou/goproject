package controllers

import (
	"encoding/json"

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
func (c *ExampleController) Examples() {
	exampleService := services.ExampleService{}
	result := exampleService.GetExampleSettings()
	c.Ctx.JSONResp(result)
}

// @Title create examples
// @Success 200 {object} 	services.ExampleModel[]
// @Param	body		body 	services.ExampleModel	true		"body for create example content"
// @router / [post]
func (c *ExampleController) CreateExample() {
	exampleService := services.ExampleService{}
	var request services.ExampleModel
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	exampleService.CreateCustomExample(request)
	c.Ctx.WriteString("ok")
}

// @Title get examples
// @Success 200 {string}
// @router /init [post]
func (c *ExampleController) InitExamples() {
	exampleService := services.ExampleService{}
	exampleService.InitExamples()
	c.Ctx.WriteString("ok")
}
