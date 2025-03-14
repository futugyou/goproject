package openai

import (
	"context"
	"os"
	"strings"

	"github.com/futugyou/yomawari/common/errorutils"

	"golang.org/x/exp/slices"
)

const createImagesPath string = "images/generations"
const editImagesPath string = "images/edits"
const variationImagesPath string = "images/variations"

var supportedImageType = []string{"png"}

type CreateImagesRequest struct {
	Prompt         string          `json:"prompt"`
	N              int32           `json:"n,omitempty"`
	Size           ImageSize       `json:"size,omitempty"`
	ResponseFormat ImageFormatType `json:"response_format,omitempty"`
	User           string          `json:"user,omitempty"`
}

type CreateImagesResponse struct {
	Error   *errorutils.OpenaiError `json:"error,omitempty"`
	Created int                     `json:"created,omitempty"`
	Data    []data                  `json:"data,omitempty"`
}

type data struct {
	URL     string `json:"url,omitempty"`
	B64Json string `json:"b64_json,omitempty"`
}

type EditImagesRequest struct {
	Image          *os.File        `json:"image"`
	Mask           *os.File        `json:"mask,omitempty"`
	Prompt         string          `json:"prompt"`
	N              int32           `json:"n,omitempty"`
	Size           ImageSize       `json:"size,omitempty"` //'256x256', '512x512', '1024x1024'
	ResponseFormat ImageFormatType `json:"response_format,omitempty"`
	User           string          `json:"user,omitempty"`
}

type EditImagesResponse struct {
	Error   *errorutils.OpenaiError `json:"error,omitempty"`
	Created int                     `json:"created,omitempty"`
	Data    []data                  `json:"data,omitempty"`
}

type VariationImagesRequest struct {
	Image          *os.File        `json:"image"`
	N              int32           `json:"n,omitempty"`
	Size           ImageSize       `json:"size,omitempty"` //'256x256', '512x512', '1024x1024'
	ResponseFormat ImageFormatType `json:"response_format,omitempty"`
	User           string          `json:"user,omitempty"`
}

type VariationImagesResponse struct {
	Error   *errorutils.OpenaiError `json:"error,omitempty"`
	Created int                     `json:"created,omitempty"`
	Data    []data                  `json:"data,omitempty"`
}

type ImageService service

func (c *ImageService) CreateImages(ctx context.Context, request CreateImagesRequest) *CreateImagesResponse {
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

	c.client.httpClient.Post(ctx, createImagesPath, request, result)
	return result
}

func (c *ImageService) EditImages(ctx context.Context, request EditImagesRequest) *EditImagesResponse {
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
		result.Error = errorutils.MessageError("Images can nod be nil.")
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

	c.client.httpClient.PostWithFile(ctx, editImagesPath, &request, result)
	return result
}

func (c *ImageService) VariationImages(ctx context.Context, request VariationImagesRequest) *VariationImagesResponse {
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
		result.Error = errorutils.MessageError("Images can nod be nil.")
		return result
	}

	err = validateImageType(request.Image)
	if err != nil {
		result.Error = err
		return result
	}

	c.client.httpClient.PostWithFile(ctx, variationImagesPath, &request, result)
	return result
}

func validateImageSize(size ImageSize) *errorutils.OpenaiError {
	if len(size) == 0 || !slices.Contains(SupportededImageSize, size) {
		return errorutils.UnsupportedTypeError("images size", size, SupportededImageSize)
	}

	return nil
}

func validateImageResponseFormat(format ImageFormatType) *errorutils.OpenaiError {
	if len(format) == 0 || !slices.Contains(SupportedImageResponseFormat, format) {
		return errorutils.UnsupportedTypeError("ResponseFormat", format, SupportedImageResponseFormat)
	}

	return nil
}

func validateImageType(file *os.File) *errorutils.OpenaiError {
	if file == nil {
		return nil
	}

	segmentations := strings.Split(file.Name(), ".")
	if len(segmentations) <= 1 {
		return errorutils.UnsupportedTypeError("Image type", "nil", supportedImageType)
	}

	suffix := strings.ToLower(strings.Split(file.Name(), ".")[len(segmentations)-1])
	if !slices.Contains(supportedAudioType, suffix) {
		return errorutils.UnsupportedTypeError("Image type", suffix, supportedImageType)
	}

	return nil
}
