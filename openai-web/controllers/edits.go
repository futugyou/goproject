package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/server/web"

	"github.com/futugyousuzu/go-openai-web/services"
)

// Operations about edit
type EditController struct {
	web.Controller
}

// @Title Create edit
// @Description create edit
// @Param	body		body 	services.CreateEditsRequest	true		"body for create edit content"
// @Success 200 {test} 	string
// @router / [post]
func (c *EditController) CreateEdit() {
	var r services.CreateEditsRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &r)

	completionService := services.NewEditService(createOpenAICLient())
	result := completionService.CreateEdit(r)
	c.Ctx.JSONResp(result)
}
