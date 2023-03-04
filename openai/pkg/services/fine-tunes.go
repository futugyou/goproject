package services

import (
	"openai/lib"

	"github.com/beego/beego/v2/core/config"
)

type FineTuneService struct {
}

func (s *FineTuneService) ListFineTuneEventsStream(fine_tune_id string) *lib.ListFinetuneEventResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	result := client.ListFinetuneEventsStream(fine_tune_id)
	return result
}
