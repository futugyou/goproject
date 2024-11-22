package controllers

import (
	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Models
type ModelController struct {
	web.Controller
}

// @Title ListModel
// @Description list model
// @Success 200 {list model} 	[]string
// @router / [get]
func (c *ModelController) ListModel() {
	modelService := services.NewModelService(createOpenAICLient(), createRedisICLient())
	result := modelService.GetAllModels(c.Ctx.Request.Context())
	c.Ctx.JSONResp(result)
}
