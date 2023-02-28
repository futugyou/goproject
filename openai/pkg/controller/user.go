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
	// result := services.CreateEmbeddings()
	// result := services.ListFiles()
	// result := services.UploadFiles()
	// u.Ctx.WriteString(result)

	// result := services.ListModelsLib()
	// result := services.RetrieveModelLib()
	// result := services.CreateCompletionLib()
	// result := services.CreateEditsLib()
	// result := services.CreateImagesLib()
	// result := services.EditImageslib()
	// result := services.VariationImagesLib()
	// result := services.CreateEmbeddingslib()
	// result := services.ListFilesLib()
	result := services.UploadFileslib()
	u.Ctx.JSONResp(result)

}
