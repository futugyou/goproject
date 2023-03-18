package openai
 
type openaiClient struct {
	httpClient IHttpClient
}

func NewClient(apikey string) *openaiClient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}

	return &openaiClient{
		httpClient: NewHttpClient(apikey),
	}
}

func (c *openaiClient) SetOrganization(organization string) {
	c.httpClient.SetOrganization(organization)
}

func (c *openaiClient) SetBaseUrl(baseurl string) {
	c.httpClient.SetBaseUrl(baseurl)
}
