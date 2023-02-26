package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/beego/beego/v2/core/config"

	lib "openai/lib/pkg"
)

type Payload struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int64   `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int64   `json:"top_p"`
	FrequencyPenalty int64   `json:"frequency_penalty"`
	PresencePenalty  int64   `json:"presence_penalty"`
	Model            string  `json:"model"`
}

func Completions() string {

	data := Payload{
		Prompt:           "how to use github",
		MaxTokens:        2048,
		Temperature:      0.5,
		TopP:             0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "text-davinci-003",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey, _ := config.String("openaikey")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", openaikey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(all)
}

func CreateCompletionLib() interface{} {
	request := lib.CreateCompletionRequest{
		Prompt:           "how to use github",
		MaxTokens:        2048,
		Temperature:      0.5,
		Top_p:            0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "text-davinci-003",
		Logprobs:         1,
	}

	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	return client.CreateCompletion(request)
}

func ListModels() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey, _ := config.String("openaikey")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", openaikey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(all)
}

func CallLib() interface{} {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	return client.Listmodels()
}

func RetrieveModel() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models/text-davinci-003", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey, _ := config.String("openaikey")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", openaikey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(all)
}
func RetrieveModelLib() interface{} {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	return client.RetrieveModel("text-davinci-003")
}
