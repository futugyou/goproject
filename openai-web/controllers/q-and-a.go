package controllers

import (
	"encoding/json"

	"github.com/futugyousuzu/go-openai-web/models"

	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"

	"github.com/devfeel/mapper"
	"github.com/futugyousuzu/go-openai-web/services"
)

// Operations about q&a
type QuestionController struct {
	web.Controller
}

// @Title Createq&a
// @Description create q&a
// @Param	body		body 	models.QuestionAnswer	true		"body for create q&a content"
// @Success 200 {test} 	string
// @router / [post]
func (c *QuestionController) CreateQAndA() {
	var r models.QuestionAnswer
	json.Unmarshal(c.Ctx.Input.RequestBody, &r)
	valid := validation.Validation{}
	b, err := valid.Valid(&r)

	if err != nil {
		c.Ctx.JSONResp(err)
		return
	}

	if !b {
		for _, err := range valid.Errors {
			c.Ctx.WriteString(err.Key + " " + err.Message)
			return
		}
	}

	completionService := services.CompletionService{}
	co := services.CompletionModel{}
	mapper.AutoMapper(&r, &co)

	re := services.CreateCompletionRequest{
		CompletionModel: co,
	}

	result := completionService.CreateCompletion(re)
	c.Ctx.JSONResp(result)
}
