package controllers

import (
	"context"
	"encoding/json"
	"os"

	lib "github.com/futugyousuzu/go-openai"
	"github.com/futugyousuzu/go-openai-web/services"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Chat
type AudioController struct {
	web.Controller
}

// @Title CreateAudioTranscription
// @Description create audio transcription
// @Param	body		body 	lib.CreateAudioTranscriptionRequest	true		"body for create audio transcription content"
// @Success 200 {object} 	lib.CreateAudioTranscriptionResponse
// @router /transcription [post]
func (c *AudioController) CreateAudioTranscription(request lib.CreateAudioTranscriptionRequest) {
	chatService := services.NewAudioService(createOpenAICLient())
	var audio lib.CreateAudioTranscriptionRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result := chatService.CreateAudioTranscription(audio)
	c.Ctx.JSONResp(result)
}

// @Title CreateAudioTranslation
// @Description create audio translation
// @Param	body		body 	lib.CreateAudioTranslationRequest	true		"body for create audio translation content"
// @Success 200 {object} 	lib.CreateAudioTranslationResponse
// @router /translation [post]
func (c *AudioController) CreateAudioTranslation(request lib.CreateAudioTranslationRequest) {
	chatService := services.NewAudioService(createOpenAICLient())
	var audio lib.CreateAudioTranslationRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result := chatService.CreateAudioTranslation(audio)
	c.Ctx.JSONResp(result)
}

func createOpenAICLient() *lib.OpenaiClient {
	openaikey := os.Getenv("openaikey")
	openaiurl := os.Getenv("openaiurl")
	client := lib.NewClient(openaikey)
	if len(openaiurl) > 0 {
		client.SetBaseUrl(openaiurl)
	}

	return client
}

func createRedisICLient() *redis.Client {
	client, _ := services.RedisClient(os.Getenv("REDIS_URL"))
	return client
}

func createMongoDbCLient() *mongo.Client {
	uri := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	return client
}
