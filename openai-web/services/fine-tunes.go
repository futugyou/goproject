package services

import (
	"os"

	lib "github.com/futugyousuzu/go-openai"
)

type FineTuneService struct {
}

func (s *FineTuneService) ListFinetuneEvents(fine_tune_id string) *lib.ListFinetuneEventResponse {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	result := client.ListFinetuneEvents(fine_tune_id)
	return result
}
