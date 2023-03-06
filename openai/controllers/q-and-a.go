package controllers

import (
	"encoding/json"
	"openai/models"

	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
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
func (c *QuestionController) CreateQAndA(request models.QuestionAnswer) {
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

	c.Ctx.WriteString("ok")
}
