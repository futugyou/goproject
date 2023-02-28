package lib

import "os"

const createImagesPath string = "images/generations"
const editImagesPath string = "images/edits"
const variationImagesPath string = "images/variations"

type CreateImagesRequest struct {
	Prompt         string `json:"prompt"`
	N              int32  `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type CreateImagesResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Created int          `json:"created,omitempty"`
	Data    []data       `json:"data,omitempty"`
}

type data struct {
	URL     string `json:"url,omitempty"`
	B64Json string `json:"b64_json,omitempty"`
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
	Created int          `json:"created,omitempty"`
	Data    []data       `json:"data,omitempty"`
}

type VariationImagesRequest struct {
	Image          *os.File `json:"image"`
	N              int32    `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"` //'256x256', '512x512', '1024x1024'
	ResponseFormat string   `json:"response_format,omitempty"`
	User           string   `json:"user,omitempty"`
}

type VariationImagesResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Created int          `json:"created,omitempty"`
	Data    []data       `json:"data,omitempty"`
}

func (client *openaiClient) CreateImages(request CreateImagesRequest) *CreateImagesResponse {
	result := &CreateImagesResponse{}
	client.Post(createImagesPath, request, result)
	return result
}

func (client *openaiClient) EditImages(request EditImagesRequest) *EditImagesResponse {
	result := &EditImagesResponse{}
	client.PostWithFile(editImagesPath, &request, result)
	return result
}

func (client *openaiClient) VariationImages(request VariationImagesRequest) *VariationImagesResponse {
	result := &VariationImagesResponse{}
	client.PostWithFile(variationImagesPath, &request, result)
	return result
}
