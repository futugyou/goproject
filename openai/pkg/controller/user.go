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
	// result := services.RetrieveFile()
	// result := services.RetrieveFileContent()
	// result := services.DeleteFile()
	// result := services.CreateFinetune()
	// result := services.CancelFinetune()
	// result := services.ListFinetunes()
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
	// result := services.UploadFileslib()
	// result := services.RetrieveFileLib()
	// result := services.DeleteFileLib()
	// result := services.CreateFinetunelib()
	// result := services.CancelFinetunelib()
	result := services.ListFinetunesLib()
	u.Ctx.JSONResp(result)

}
