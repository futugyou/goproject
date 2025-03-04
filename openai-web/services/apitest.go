package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	lib "github.com/futugyou/ai-extension/openai"
)

type Payload struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int64   `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int64   `json:"top_p"`
	FrequencyPenalty int64   `json:"frequency_penalty"`
	PresencePenalty  int64   `json:"presence_penalty"`
	Model            string  `json:"model"`
	Stream           bool    `json:"stream"`
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
		Stream:           true,
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
	openaikey := os.Getenv("openaikey")
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

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Completion.CreateCompletion(context.Background(), request)
}

func ListModels() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func ListModelsLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Model.ListModels(context.Background())
}

func RetrieveModel() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models/text-davinci-003", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Model.RetrieveModel(context.Background(), "text-davinci-003")
}

type CreateEditsModel struct {
	Model       string `json:"model"`
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
}

func CreateEdits() string {
	data := CreateEditsModel{
		Model:       "text-davinci-edit-001",
		Input:       "What day of the wek is it?",
		Instruction: "Fix the spelling mistakes",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/edits", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CreateEditsLib() interface{} {
	request := lib.CreateEditsRequest{
		Model:       "text-davinci-edit-001",
		Input:       "What day of the wek is it?",
		Instruction: "Fix the spelling mistakes",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Edit.CreateEdits(context.Background(), request)
}

func CreateImages() string {
	data := lib.CreateImagesRequest{
		Prompt:         "A cute baby sea otter",
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "b64_json",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CreateImagesLib() interface{} {
	data := lib.CreateImagesRequest{
		Prompt: "A cute baby sea otter",
		N:      1,
		Size:   "1024x1024",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Image.CreateImages(context.Background(), data)
}

func EditImages() string {
	// image, _ := os.Create("./otter.png")
	// mask, _ := os.Create("./mask.png")
	image, _ := os.Open("./otter.png")
	mask, _ := os.Open("./mask.png")
	defer func() {
		mask.Close()
		image.Close()
		// os.Remove("mask.png")
		// os.Remove("otter.png")
	}()

	data := lib.EditImagesRequest{
		Image:  image,
		Mask:   mask,
		Prompt: "A cute baby sea otter",
		N:      1,
		Size:   "512x512",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	wimage, _ := writer.CreateFormFile("image", data.Image.Name())
	io.Copy(wimage, data.Image)
	wmask, _ := writer.CreateFormFile("image", data.Mask.Name())
	io.Copy(wmask, data.Mask)
	writer.WriteField("n", strconv.FormatInt(int64(data.N), 10))
	writer.WriteField("prompt", data.Prompt)
	writer.WriteField("size", fmt.Sprintf("%v", data.Size))

	writer.Close()
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/edits", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	openaikey := os.Getenv("openaikey")
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

func EditImageslib() interface{} {
	// image, _ := os.Create("./otter.png")
	// mask, _ := os.Create("./mask.png")
	image, _ := os.Open("./otter.png")
	// mask, _ := os.Open("./mask.png")
	defer func() {
		// mask.Close()
		image.Close()
		// os.Remove("mask.png")
		// os.Remove("otter.png")
	}()

	data := lib.EditImagesRequest{
		Image: image,
		// Mask:   mask,
		Prompt: "A cute baby sea otter",
		// N:      1,
		Size: "512x512",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Image.EditImages(context.Background(), data)
}

func VariationImagesLib() interface{} {
	image, _ := os.Open("./otter.png")
	defer func() {
		image.Close()
	}()

	data := lib.VariationImagesRequest{
		Image: image,
		Size:  "512x512",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Image.VariationImages(context.Background(), data)
}

func CreateEmbeddings() string {
	data := lib.CreateEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"The food was delicious and the waiter..."},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CreateEmbeddingslib() interface{} {
	data := lib.CreateEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"The food was delicious and the waiter..."},
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Embedding.CreateEmbeddings(context.Background(), data)
}

func ListFiles() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/files", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func ListFilesLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.File.ListFiles(context.Background())
}

func UploadFiles() string {
	file, _ := os.Open("./files.jsonl")
	defer func() {
		file.Close()
	}()

	data := lib.UploadFilesRequest{
		File:    file,
		Purpose: "fine-tune",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	wimage, _ := writer.CreateFormFile("file", data.File.Name())
	io.Copy(wimage, data.File)
	writer.WriteField("purpose", data.Purpose)

	writer.Close()
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/files", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	openaikey := os.Getenv("openaikey")
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

func UploadFileslib() interface{} {
	file, _ := os.Open("./files.jsonl")
	defer func() {
		file.Close()
	}()

	data := lib.UploadFilesRequest{
		File:    file,
		Purpose: "fine-tune",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.File.UploadFiles(context.Background(), data)
}

func RetrieveFile() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/files/file-shJO2TBQNSrDFVCXY0RNLSC2", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func RetrieveFileLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.File.RetrieveFile(context.Background(), "file-Be1Itkt0E2SinfiOnxYRPjVx")
}

func RetrieveFileContent() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/files/file-shJO2TBQNSrDFVCXY0RNLSC2/content", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func DeleteFile() string {
	req, err := http.NewRequest("DELETE", "https://api.openai.com/v1/files/file-shJO2TBQNSrDFVCXY0RNLSC2", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func DeleteFileLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.File.DeleteFile(context.Background(), "file-Be1Itkt0E2SinfiOnxYRPjVx")
}

func CreateFinetune() string {
	data := lib.CreateFinetuneRequest{
		TrainingFile:   "file-YUco6HX1ikrEK9CCUnVfDCLs",
		ValidationFile: "file-NXWeVnozaOT7ckA5gUtuVvhJ",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/fine-tunes", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CreateFinetunelib() interface{} {
	data := lib.CreateFinetuneRequest{
		TrainingFile:   "file-YUco6HX1ikrEK9CCUnVfDCLs",
		ValidationFile: "file-NXWeVnozaOT7ckA5gUtuVvhJ",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.CreateFinetune(context.Background(), data)
}

func CancelFinetune() string {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/fine-tunes/ft-c0pBGCqr0daPhapyJgJXxHJp/cancel", nil)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CancelFinetunelib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.CancelFinetune(context.Background(), "ft-wVjb6K7ngTeYeW6QT1eDQikZ")
}

func ListFinetunes() string {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/fine-tunes", nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func ListFinetunesLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.ListFinetune(context.Background())
}

func RetrieveFinetunelib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.RetrieveFinetune(context.Background(), "ft-W0GCdkAnSKNIoyWhfbe86zzv")
}

func ListFinetuneEventslib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.ListFinetuneEvents(context.Background(), "ft-W0GCdkAnSKNIoyWhfbe86zzv")
}

func DeleteFinetuneMdelLib() interface{} {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Finetune.DeleteFinetuneMdel(context.Background(), "curie:ft-personal-2023-02-28-05-52-07")
}

func CreateModeration() string {
	data := lib.CreateModerationRequest{
		Input: "how to use github",
		Model: "text-moderation-latest",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/moderations", body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	openaikey := os.Getenv("openaikey")
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

func CreateModerationLib() interface{} {
	request := lib.CreateModerationRequest{
		Input: "how to use github",
		Model: "text-moderation-latest",
	}

	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Moderation.CreateModeration(context.Background(), request)
}

func CreateAudioTranscriptionLib() interface{} {
	file, _ := os.Open("./multilingual.mp3")
	defer func() {
		file.Close()
	}()

	data := lib.CreateAudioTranscriptionRequest{
		File:           file,
		Model:          "whisper-1",
		Prompt:         "this is test",
		ResponseFormat: "json",
		Temperature:    0.5,
		Language:       "en",
	}
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	return client.Audio.CreateAudioTranscription(context.Background(), data)
}

func CreateAudioTranscription() string {
	file, _ := os.Open("./multilingual.mp3")
	defer func() {
		file.Close()
	}()

	data := lib.CreateAudioTranscriptionRequest{
		File:           file,
		Model:          "whisper-1",
		Prompt:         "this is test",
		ResponseFormat: "json",
		Temperature:    0.5,
		Language:       "en",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	wimage, _ := writer.CreateFormFile("file", data.File.Name())
	io.Copy(wimage, data.File)
	writer.WriteField("model", data.Model)
	writer.WriteField("prompt", data.Prompt)
	writer.WriteField("response_format", fmt.Sprintf("%v", data.ResponseFormat))
	writer.WriteField("temperature", fmt.Sprintf("%f", data.Temperature))
	writer.WriteField("language", data.Language)

	writer.Close()
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		log.Println(err.Error())

	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	openaikey := os.Getenv("openaikey")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", openaikey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())

	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(all)
}
