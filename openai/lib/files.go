package lib

import "os"

const listFilesPath string = "files"
const uploadFilesPath string = "files"

type UploadFilesRequest struct {
	File    *os.File `json:"file"`
	Purpose string   `json:"purpose"`
}

type UploadFilesResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	fileModel
}

type ListFilesResponse struct {
	Error  *OpenaiError `json:"error,omitempty"`
	Object string       `json:"object,omitempty"`
	Data   []fileModel  `json:"data,omitempty"`
}

type fileModel struct {
	Object        string      `json:"object"`
	ID            string      `json:"id"`
	Purpose       string      `json:"purpose"`
	Filename      string      `json:"filename"`
	Bytes         int         `json:"bytes"`
	CreatedAt     int         `json:"created_at"`
	Status        string      `json:"status"`
	StatusDetails interface{} `json:"status_details"`
}

func (client *openaiClient) ListFiles() *ListFilesResponse {
	result := &ListFilesResponse{}
	client.Get(listFilesPath, result)
	return result
}
