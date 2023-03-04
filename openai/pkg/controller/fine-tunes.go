package controller

import (
	"openai/pkg/services"

	"github.com/beego/beego/v2/server/web"
)

type FineTuneController struct {
	web.Controller
}

func (c *FineTuneController) ListFineTuneEvent() {
	fine_tune_id := c.GetString("fine_tune_id")
	modelService := services.FineTuneService{}
	result := modelService.ListFineTuneEventsStream(fine_tune_id)
	c.Ctx.JSONResp(result)
}
