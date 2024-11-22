package openai

import (
	"fmt"
	"os"
)

const listFilesPath string = "files"
const uploadFilesPath string = "files"
const retrieveFilePath string = "files/%s"
const retrieveFileContentPath string = "files/%s/content"
const deleteFilePath string = "files/%s"

type UploadFilesRequest struct {
	File    *os.File `json:"file"`
	Purpose string   `json:"purpose"`
}

type UploadFilesResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FileModel
}

type ListFilesResponse struct {
	Error  *OpenaiError `json:"error,omitempty"`
	Object string       `json:"object,omitempty"`
	Data   []FileModel  `json:"data,omitempty"`
}

type RetrieveFileResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FileModel
}

type RetrieveFileContentResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	FileModel
}

type DeleteFileResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Object  string       `json:"object,omitempty"`
	ID      string       `json:"id,omitempty"`
	Deleted bool         `json:"deleted,omitempty"`
}

type FileService service

func (c *FileService) ListFiles() *ListFilesResponse {
	result := &ListFilesResponse{}
	c.client.httpClient.Get(listFilesPath, result)
	return result
}

func (c *FileService) UploadFiles(request UploadFilesRequest) *UploadFilesResponse {
	result := &UploadFilesResponse{}
	c.client.httpClient.PostWithFile(uploadFilesPath, &request, result)
	return result
}

func (c *FileService) RetrieveFile(file_id string) *RetrieveFileResponse {
	result := &RetrieveFileResponse{}
	c.client.httpClient.Get(fmt.Sprintf(retrieveFilePath, file_id), result)
	return result
}

// To help mitigate abuse, downloading of fine-tune training files is disabled for free accounts
// so i do not know the response type
// And saw some other lib, and i think it is a FileModel type.
func (c *FileService) RetrieveFileContent(file_id string) *RetrieveFileContentResponse {
	result := &RetrieveFileContentResponse{}
	c.client.httpClient.Get(fmt.Sprintf(retrieveFileContentPath, file_id), result)
	return result
}

func (c *FileService) DeleteFile(file_id string) *DeleteFileResponse {
	result := &DeleteFileResponse{}
	c.client.httpClient.Delete(fmt.Sprintf(deleteFilePath, file_id), result)
	return result
}
