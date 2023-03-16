package lib

import (
	c "openai/lib/internal"
)

type openaiClient struct {
	httpClient c.IHttpClient
}

func NewClient(apikey string) *openaiClient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}

	return &openaiClient{
		httpClient: c.NewHttpClient(apikey),
	}
}

func (c *openaiClient) SetOrganization(organization string) {
	c.httpClient.SetOrganization(organization)
}

func (c *openaiClient) SetBaseUrl(baseurl string) {
	c.httpClient.SetBaseUrl(baseurl)
}
