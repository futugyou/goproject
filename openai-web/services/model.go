package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	lib "github.com/futugyousuzu/go-openai"
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

	if len(rmap) > 0 {
		for _, r := range rmap {
			m := lib.Model{}
			json.Unmarshal([]byte(r), &m)
			result = append(result, ModelListResponse{Name: m.ID})
		}

		return result
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	models := client.ListModels()
	rset := make(map[string]interface{})
	if len(models.Datas) > 0 {
		for _, model := range models.Datas {
			result = append(result, ModelListResponse{Name: model.ID})
			modelstring, _ := json.Marshal(model)
			rset[model.ID] = string(modelstring)
		}
	}

	count, err := Rbd.HSet(ctx, GetAllModelsKey, rset).Result()
	if err != nil {
		logs.Error(err)
	} else {
		logs.Info(fmt.Sprintf("data count: %d", count))
	}

	return result
}

type MyHash struct {
	Key1 string `redis:"key1"`
	Key2 int    `redis:"key2"`
}
