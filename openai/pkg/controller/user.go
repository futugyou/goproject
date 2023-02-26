package controller

import (
	"openai/pkg/services"

	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

func (u *UserController) HelloWorld() {
	// u.Ctx.WriteString("hello, world")
	// result := services.Completions()
	// result := services.ListModels()
	// result := services.RetrieveModel()
	// u.Ctx.WriteString(result)
	// result := services.CallLib()
	// u.Ctx.JSONResp(result)
	// result := services.RetrieveModelLib()
	result := services.CreateCompletionLib()
	u.Ctx.JSONResp(result)

}
