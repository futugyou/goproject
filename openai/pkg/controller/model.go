package controller

import (
	"openai/pkg/services"

	"github.com/beego/beego/v2/server/web"
)

type ModelController struct {
	web.Controller
}

func (c *ModelController) ListModel() {
	modelService := services.ModelService{}
	result := modelService.GetAllModels()
	c.Ctx.JSONResp(result)
}
