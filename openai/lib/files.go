package lib

import (
	"fmt"
	"os"
)

const listFilesPath string = "files"
const uploadFilesPath string = "files"
const retrieveFilePath string = "files/%s"
const retrieveFileContentPath string = "files/%s/content"

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

type RetrieveFileResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	fileModel
}

type RetrieveFileContentResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	content
}

type content interface{}

func (client *openaiClient) ListFiles() *ListFilesResponse {
	result := &ListFilesResponse{}
	client.Get(listFilesPath, result)
	return result
}

func (client *openaiClient) UploadFiles(request UploadFilesRequest) *UploadFilesResponse {
	result := &UploadFilesResponse{}
	client.PostWithFile(uploadFilesPath, &request, result)
	return result
}

func (client *openaiClient) RetrieveFile(file_id string) *RetrieveFileResponse {
	result := &RetrieveFileResponse{}
	client.Get(fmt.Sprintf(retrieveFilePath, file_id), result)
	return result
}

// To help mitigate abuse, downloading of fine-tune training files is disabled for free accounts
// so i do not know the response type
func (client *openaiClient) RetrieveFileContent(file_id string) *RetrieveFileContentResponse {
	result := &RetrieveFileContentResponse{}
	client.Get(fmt.Sprintf(retrieveFileContentPath, file_id), result)
	return result
}
