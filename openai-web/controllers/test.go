package controllers

import (
	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about test
type TestController struct {
	web.Controller
}

// @Title test openai lib
// @Description do the test
// @router / [get]
func (c *TestController) Test() {
	result := services.CreateCompletionLib()
	c.Ctx.JSONResp(result)
}
