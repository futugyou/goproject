package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/futugyousuzu/go-openai-web/services"

	"github.com/futugyousuzu/go-openai-web/models"

	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"

	"github.com/devfeel/mapper"
)

// Operations about completion
type CompletionController struct {
	web.Controller
}

// @Title Create Completion With SSE
// @Description create completion stream
// @Param	body		body 	lib.CreateCompletionRequest	true		"body for create completion content"
// @Success 200 {object} 	lib.CreateCompletionResponse
// @router /sse [post]
func (c *CompletionController) CreateCompletionWithSSE() {
	var r models.CompletionModel
	json.Unmarshal(c.Ctx.Input.RequestBody, &r)
	valid := validation.Validation{}
	b, err := valid.Valid(&r)

	if err != nil {
		c.Ctx.JSONResp(err)
		return
	}

	if !b {
		errString := ""
		for _, err := range valid.Errors {
			errString += (err.Key + " " + err.Message)
		}

		c.Ctx.WriteString(errString)
		return
	}

	completionService := services.CompletionService{}
	co := services.CompletionModel{}
	mapper.AutoMapper(&r, &co)

	re := services.CreateCompletionRequest{
		CompletionModel: co,
	}

	result := completionService.CreateCompletionSSE(re)

	var rw = c.Ctx.ResponseWriter
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set(`Content-Type`, `text/event-stream;charset-utf-8`)
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	for response := range result {
		message := strings.Join(response.Texts, ",")
		if len(message) == 0 {
			continue
		}

		message = url.QueryEscape(message)
		data := fmt.Sprintf("data: %s\n\n", message)
		rw.Write([]byte(data))
		rw.Flush()
	}
	rw.Write([]byte("data: [DONE]\n\n"))
	rw.Flush()
}

// @Title Create Completion
// @Description create completion
// @Param	body		body 	models.CompletionModel	true		"body for create completion content"
// @router / [post]
func (c *CompletionController) CreateCompletion() {
	var r models.CompletionModel
	json.Unmarshal(c.Ctx.Input.RequestBody, &r)
	valid := validation.Validation{}
	b, err := valid.Valid(&r)

	if err != nil {
		c.Ctx.JSONResp(err)
		return
	}

	if !b {
		errString := ""
		for _, err := range valid.Errors {
			errString += (err.Key + " " + err.Message)
		}

		c.Ctx.WriteString(errString)
		return
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
