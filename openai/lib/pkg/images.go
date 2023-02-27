package pkg

const createImagesPath string = "images/generations"

type CreateImagesRequest struct {
	Prompt         string `json:"prompt"`
	N              int32  `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}
