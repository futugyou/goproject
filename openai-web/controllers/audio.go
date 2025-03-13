package controllers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/futugyousuzu/go-openai-web/services"
	lib "github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
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
// @Param	body		body 	lib.AudioTranscriptionNewParams	true		"body for create audio transcription content"
// @Success 200 {object} 	lib.CreateAudioTranscriptionResponse
// @router /transcription [post]
func (c *AudioController) CreateAudioTranscription(request lib.AudioTranscriptionNewParams) {
	chatService := services.NewAudioService(createOpenAICLient())
	var audio lib.AudioTranscriptionNewParams
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result, err := chatService.CreateAudioTranscription(c.Ctx.Request.Context(), audio)
	if err != nil {
		c.Ctx.JSONResp(err)
		return
	}
	c.Ctx.JSONResp(result)
}

// @Title CreateAudioTranslation
// @Description create audio translation
// @Param	body		body 	lib.AudioTranslationNewParams	true		"body for create audio translation content"
// @Success 200 {object} 	lib.CreateAudioTranslationResponse
// @router /translation [post]
func (c *AudioController) CreateAudioTranslation(request lib.AudioTranslationNewParams) {
	chatService := services.NewAudioService(createOpenAICLient())
	var audio lib.AudioTranslationNewParams
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result, err := chatService.CreateAudioTranslation(c.Ctx.Request.Context(), audio)
	if err != nil {
		c.Ctx.JSONResp(err)
		return
	}
	c.Ctx.JSONResp(result)
}

func createOpenAICLient() *lib.Client {
	openaikey := os.Getenv("openaikey")
	openaiurl := os.Getenv("openaiurl")

	const azureOpenAIAPIVersion = "2024-06-01"

	client := lib.NewClient(
		azure.WithEndpoint(openaiurl, azureOpenAIAPIVersion),
		azure.WithAPIKey(openaikey),
	)

	return client
}

func createRedisICLient() *redis.Client {
	client, _ := services.RedisClient(os.Getenv("REDIS_URL"))
	return client
}

func createMongoDbCLient(ctx context.Context) *mongo.Client {
	uri := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client
}
