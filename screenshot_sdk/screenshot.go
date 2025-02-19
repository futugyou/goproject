package screenshot

import "net/http"

type ScreenshotClient struct {
	http              *http.Client
	key               string
	common            service
	Apiflash          *ApiflashService
	Screenshotmachine *ScreenshotmachineService
}

type service struct {
	client *ScreenshotClient
}

func NewClient(token string) *ScreenshotClient {
	c := &ScreenshotClient{
		http: &http.Client{},
		key:  token,
	}
	c.initialize()
	return c
}

func (c *ScreenshotClient) initialize() {
	c.common.client = c
	c.Apiflash = (*ApiflashService)(&c.common)
	c.Screenshotmachine = (*ScreenshotmachineService)(&c.common)
}

type BaseResponse struct {
	Message *string `json:"message,omitempty"`
}
