package gofile

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type ContentService service

func (s *ContentService) UploadFile(ctx context.Context, request UploadFileRequest) (*UploadFileResponse, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	filePart, err := writer.CreateFormFile("file", request.FileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(filePart, request.File)
	if err != nil {
		return nil, err
	}

	if request.FolderId != nil {
		err = writer.WriteField("folderId", *request.FolderId)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	server := request.Server
	if server == nil {
		serverinfo, err := s.client.Servers.GetServer(ctx)
		if err != nil {
			return nil, err
		}

		if len(serverinfo.Data.Servers) == 0 {
			return nil, fmt.Errorf("no server available")
		}

		server = &serverinfo.Data.Servers[0].Name
	}

	url := fmt.Sprintf("https://%s.gofile.io/contents/uploadfile", *server)
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+s.client.key)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload error, status code: %v", resp.StatusCode)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &UploadFileResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

type UploadFileRequest struct {
	File     io.Reader
	FolderId *string
	Server   *string
	FileName string
}

type UploadFileResponse struct {
	Data   FileData `json:"data"`
	Status string   `json:"status"`
}

type FileData struct {
	CreateTime       int64    `json:"createTime"`
	DownloadPage     string   `json:"downloadPage"`
	ID               string   `json:"id"`
	Md5              string   `json:"md5"`
	Mimetype         string   `json:"mimetype"`
	ModTime          int64    `json:"modTime"`
	Name             string   `json:"name"`
	ParentFolder     string   `json:"parentFolder"`
	ParentFolderCode string   `json:"parentFolderCode"`
	Servers          []string `json:"servers"`
	Size             int64    `json:"size"`
	Type             string   `json:"type"`
}

func (s *ContentService) CreateFolder(ctx context.Context, request CreateFolderRequest) (*CreateFolderResponse, error) {
	path := "https://api.gofile.io/contents/createFolder"
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+s.client.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create folder error, status code: %v", resp.StatusCode)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &CreateFolderResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

type CreateFolderRequest struct {
	FolderName     *string `json:"folderName"`
	ParentFolderId string  `json:"parentFolderId"`
}

type CreateFolderResponse struct {
	Status string     `json:"status"`
	Data   FolderData `json:"data"`
}

type FolderData struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ParentFolder string `json:"parentFolder"`
	CreateTime   int64  `json:"createTime"`
	ModTime      int64  `json:"modTime"`
	Code         string `json:"code"`
}

func (s *ContentService) UpdateContent(ctx context.Context, request UpdateContentRequest) (*UpdateContentResponse, error) {
	path := fmt.Sprintf("https://api.gofile.io/contents/%s/update", request.ContentId)

	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+s.client.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update content error, status code: %v", resp.StatusCode)
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &UpdateContentResponse{}
	if err = json.Unmarshal(all, response); err != nil {
		return nil, err
	}

	return response, nil
}

// name: Content name (files & folders)
//
// description: Download page description (folders only)
//
// tags: Comma-separated tags (folders only)
//
// public: Public access status (folders only)
//
// expiry: Expiration date timestamp (folders only)
//
// password: Access password (folders only)
type UpdateContentRequest struct {
	ContentId      string `json:"-"`
	Attribute      string `json:"attribute"`
	AttributeValue string `json:"attributeValue"`
}

type UpdateContentResponse struct {
	Status string      `json:"status"`
	Data   ContentInfo `json:"data"`
}

type ContentInfo struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	CreateTime   int64  `json:"createTime"`
	ModTime      int64  `json:"modTime"`
	ParentFolder string `json:"parentFolder"`
}
