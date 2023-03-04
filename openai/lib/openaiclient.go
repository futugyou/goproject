package lib

import "net/http"

const baseUrl string = "https://api.openai.com/v1/"

type openaiClient struct {
	httpClient *httpClient
}

func NewClient(apikey string) *openaiClient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}

	return &openaiClient{
		httpClient: &httpClient{
			apikey:       apikey,
			organization: "",
			baseurl:      baseUrl,
			http:         &http.Client{},
		},
	}
}

func (c *openaiClient) SetOrganization(organization string) {
	c.httpClient.organization = organization
}

func (c *openaiClient) SetBaseUrl(baseurl string) {
	c.httpClient.baseurl = baseurl
}
