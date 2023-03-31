package services

import (
	lib "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
)

type ModelListResponse struct {
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

type ModelService struct {
}

func (s *ModelService) GetAllModels() []ModelListResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	result := make([]ModelListResponse, 0)
	models := client.ListModels()
	if len(models.Datas) > 0 {
		for _, model := range models.Datas {
			result = append(result, ModelListResponse{Name: model.ID})
		}
	}

	return result
}
