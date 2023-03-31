package services

import (
	"encoding/json"

	lib "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
)

type ModelListResponse struct {
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

type ModelService struct {
}

const GetAllModelsKey string = "GetAllModelsKey"

func (s *ModelService) GetAllModels() []ModelListResponse {
	result := make([]ModelListResponse, 0)
	rmap, _ := Rbd.HGetAll(ctx, GetAllModelsKey).Result()

	if rmap != nil {
		for _, r := range rmap {
			m := lib.Model{}
			json.Unmarshal([]byte(r), &m)
			result = append(result, ModelListResponse{Name: m.ID})
		}
		return result
	}

	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	models := client.ListModels()
	rset := make(map[string]interface{})
	if len(models.Datas) > 0 {
		for _, model := range models.Datas {
			result = append(result, ModelListResponse{Name: model.ID})
			rset[model.ID] = model
		}
	}

	Rbd.HSet(ctx, GetAllModelsKey, rset)
	return result
}
