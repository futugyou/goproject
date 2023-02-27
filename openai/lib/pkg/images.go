package pkg

const createImagesPath string = "images/generations"

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
	URL string `json:"url"`
}

func (client *openaiClient) CreateImages(request CreateImagesRequest) *CreateImagesResponse {
	result := &CreateImagesResponse{}
	client.Post(createImagesPath, request, result)
	return result
}
