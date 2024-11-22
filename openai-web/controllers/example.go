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

// @Title get all examples
// @Success 200 {object} 	services.ExampleModel
// @router / [get]
func (c *ExampleController) Examples() {
	typestring := c.GetString("type")

	exampleService := services.NewExampleService(createMongoDbCLient(c.Ctx.Request.Context()), createRedisICLient())
	var result []services.ExampleModel
	if typestring == "custom" {
		result = exampleService.GetCustomExamples(c.Ctx.Request.Context())
	} else {
		result = exampleService.GetSystemExamples(c.Ctx.Request.Context())
	}
	c.Ctx.JSONResp(result)
}

// @Title create examples
// @Success 200 {object} 	services.ExampleModel[]
// @Param	body		body 	services.ExampleModel	true		"body for create example content"
// @router / [post]
func (c *ExampleController) CreateExample() {
	typestring := c.GetString("type")
	exampleService := services.NewExampleService(createMongoDbCLient(c.Ctx.Request.Context()), createRedisICLient())
	var request services.ExampleModel
	json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if len(request.Key) == 0 {
		c.Ctx.WriteString("error")
		return
	}

	if typestring == "custom" {
		exampleService.CreateCustomExample(c.Ctx.Request.Context(), request)
	} else {
		exampleService.CreateSystemExample(c.Ctx.Request.Context(), request)
	}

	c.Ctx.WriteString("ok")
}

// @Title init examples
// @Success 200 {string}
// @router /init [post]
func (c *ExampleController) InitExamples() {
	exampleService := services.NewExampleService(createMongoDbCLient(c.Ctx.Request.Context()), createRedisICLient())
	exampleService.InitExamples(c.Ctx.Request.Context())
	c.Ctx.WriteString("ok")
}

// @Title set examples
// @Success 200 {string}
// @router /reset [post]
func (c *ExampleController) ResetExamples() {
	exampleService := services.NewExampleService(createMongoDbCLient(c.Ctx.Request.Context()), createRedisICLient())
	exampleService.Reset(c.Ctx.Request.Context())
	c.Ctx.WriteString("ok")
}
