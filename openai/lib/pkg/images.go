package pkg

import "os"

const createImagesPath string = "images/generations"
const editImagesPath string = "images/edit"

type CreateImagesRequest struct {
	Prompt         string `json:"prompt"`
	N              int32  `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type CreateImagesResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Created int          `json:"created"`
	Data    []data       `json:"data"`
}

type data struct {
	URL     string `json:"url"`
	B64Json string `json:"b64_json"`
}

type EditImagesRequest struct {
	Image          *os.File `json:"image"`
	Mask           *os.File `json:"mask,omitempty"`
	Prompt         string   `json:"prompt"`
	N              int32    `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"` //'256x256', '512x512', '1024x1024'
	ResponseFormat string   `json:"response_format,omitempty"`
	User           string   `json:"user,omitempty"`
}

type EditImagesResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Created int          `json:"created"`
	Data    []data       `json:"data"`
}

func (client *openaiClient) CreateImages(request CreateImagesRequest) *CreateImagesResponse {
	result := &CreateImagesResponse{}
	client.Post(createImagesPath, request, result)
	return result
}
