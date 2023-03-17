package lib

import (
	"os"
	"strings"

	e "openai/lib/internal"

	"golang.org/x/exp/slices"
)

const createImagesPath string = "images/generations"
const editImagesPath string = "images/edits"
const variationImagesPath string = "images/variations"

var supportededImageSize = []string{"256x256", "512x512", "1024x1024"}
var supportedImageResponseFormat = []string{"url", "b64_json"}
var supportedImageType = []string{"png"}

type CreateImagesRequest struct {
	Prompt         string `json:"prompt"`
	N              int32  `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type CreateImagesResponse struct {
	Error   *e.OpenaiError `json:"error,omitempty"`
	Created int            `json:"created,omitempty"`
	Data    []data         `json:"data,omitempty"`
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
	Error   *e.OpenaiError `json:"error,omitempty"`
	Created int            `json:"created,omitempty"`
	Data    []data         `json:"data,omitempty"`
}

type VariationImagesRequest struct {
	Image          *os.File `json:"image"`
	N              int32    `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"` //'256x256', '512x512', '1024x1024'
	ResponseFormat string   `json:"response_format,omitempty"`
	User           string   `json:"user,omitempty"`
}

type VariationImagesResponse struct {
	Error   *e.OpenaiError `json:"error,omitempty"`
	Created int            `json:"created,omitempty"`
	Data    []data         `json:"data,omitempty"`
}

func (c *openaiClient) CreateImages(request CreateImagesRequest) *CreateImagesResponse {
	result := &CreateImagesResponse{}

	err := validateImageSize(request.Size)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateImageResponseFormat(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.Post(createImagesPath, request, result)
	return result
}

func (c *openaiClient) EditImages(request EditImagesRequest) *EditImagesResponse {
	result := &EditImagesResponse{}

	err := validateImageSize(request.Size)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateImageResponseFormat(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	if request.Image == nil {
		result.Error = e.MessageError("Images can nod be nil.")
		return result
	}

	err = validateImageType(request.Image)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateImageType(request.Mask)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.PostWithFile(editImagesPath, &request, result)
	return result
}

func (c *openaiClient) VariationImages(request VariationImagesRequest) *VariationImagesResponse {
	result := &VariationImagesResponse{}

	err := validateImageSize(request.Size)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateImageResponseFormat(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	if request.Image == nil {
		result.Error = e.MessageError("Images can nod be nil.")
		return result
	}

	err = validateImageType(request.Image)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.PostWithFile(variationImagesPath, &request, result)
	return result
}

func validateImageSize(size string) *e.OpenaiError {
	if len(size) == 0 || !slices.Contains(supportedAudioModel, size) {
		return e.UnsupportedTypeError("images size", size, supportededImageSize)
	}

	return nil
}

func validateImageResponseFormat(format string) *e.OpenaiError {
	if len(format) == 0 || !slices.Contains(supportedAudioModel, format) {
		return e.UnsupportedTypeError("ResponseFormat", format, supportedImageResponseFormat)
	}

	return nil
}

func validateImageType(file *os.File) *e.OpenaiError {
	if file == nil {
		return nil
	}

	segmentations := strings.Split(file.Name(), ".")
	if len(segmentations) <= 1 {
		return e.UnsupportedTypeError("Image type", "nil", supportedImageType)
	}

	suffix := strings.ToLower(strings.Split(file.Name(), ".")[len(segmentations)-1])
	if !slices.Contains(supportedAudioType, suffix) {
		return e.UnsupportedTypeError("Image type", suffix, supportedImageType)
	}

	return nil
}
