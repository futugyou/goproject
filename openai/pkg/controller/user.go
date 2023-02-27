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
	// result := services.CreateEdits()
	// result := services.CreateImages()
	// result := services.EditImages()
	// u.Ctx.WriteString(result)

	// result := services.CallLib()
	// result := services.RetrieveModelLib()
	// result := services.CreateCompletionLib()
	// result := services.CreateEditsLib()
	// result := services.CreateImagesLib()
	// result := services.EditImageslib()
	result := services.VariationImagesLib()
	u.Ctx.JSONResp(result)

}
