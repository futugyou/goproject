package openai

type OpenaiClient struct {
	httpClient *HttpClient
	common     service
	Audio      *AudioService
	Chat       *ChatService
	Completion *CompletionService
	Edit       *EditService
	Embedding  *EmbeddingService
	File       *FileService
	Finetune   *FinetuneService
	Image      *ImageService
	Model      *ModelService
	Moderation *ModerationService
}

type service struct {
	client *OpenaiClient
}

func NewClient(apikey string) *OpenaiClient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}

	c := &OpenaiClient{
		httpClient: NewHttpClient(apikey),
	}

	c.initialize()
	return c
}

func (c *OpenaiClient) SetOrganization(organization string) {
	c.httpClient.SetOrganization(organization)
}

func (c *OpenaiClient) SetBaseUrl(baseurl string) {
	c.httpClient.SetBaseUrl(baseurl)
}

func (c *OpenaiClient) initialize() {
	c.common.client = c
	c.Audio = (*AudioService)(&c.common)
	c.Chat = (*ChatService)(&c.common)
	c.Completion = (*CompletionService)(&c.common)
	c.Edit = (*EditService)(&c.common)
	c.Embedding = (*EmbeddingService)(&c.common)
	c.File = (*FileService)(&c.common)
	c.Finetune = (*FinetuneService)(&c.common)
	c.Image = (*ImageService)(&c.common)
	c.Model = (*ModelService)(&c.common)
	c.Moderation = (*ModerationService)(&c.common)
}
