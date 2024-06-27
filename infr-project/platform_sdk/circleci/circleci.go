package circleci

import (
	sdk "github.com/futugyou/infr-project/platform_sdk"
)

type CircleciClient struct {
	http sdk.IHttpClient
}

const circleci_url string = "https://circleci.com/api/v2"

func NewCircleciClient(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: sdk.NewHttpClientWithHeader(circleci_url, header),
	}
	return c
}
