package controllers

import (
	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about fine tune
type FineTuneController struct {
	web.Controller
}

// @Title Get Fine Tune Events
// @Description get fine tune by fine_tune_id
// @Param	fine_tune_id		path 	string	true		"The key for fine_tune"
// @Success 200 {object} 	lib.ListFinetuneEventResponse
// @Failure 403 fine_tune_id is empty
// @router /:fine_tune_id/events [get]
func (c *FineTuneController) ListFineTuneEvent() {
	fine_tune_id := c.GetString("fine_tune_id")
	modelService := services.FineTuneService{}
	result := modelService.ListFinetuneEvents(fine_tune_id)
	c.Ctx.JSONResp(result)
}
