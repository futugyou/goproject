package vercel

import (
	sdk "github.com/futugyou/infr-project/platform_sdk"
)

type VercelClient struct {
	token string
	http  sdk.IHttpClient
}

const vercle_url string = "https://api.vercel.com"

func NewVercelClient(token string) *VercelClient {
	c := &VercelClient{
		http: sdk.NewHttpClient(token, vercle_url),
	}
	c.token = token
	return c
}
