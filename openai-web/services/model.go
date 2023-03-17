package services

import (
	lib "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
)

type ModelService struct {
}

func (s *ModelService) GetAllModels() []string {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	result := make([]string, 0)
	models := client.ListModels()
	if len(models.Datas) > 0 {
		for _, model := range models.Datas {
			result = append(result, model.ID)
		}
	}

	return result
}
