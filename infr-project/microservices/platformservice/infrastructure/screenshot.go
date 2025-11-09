package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/futugyou/gofile"
	"github.com/futugyou/screenshot"
)
type Screenshot struct {
	screenshotClient *screenshot.ScreenshotClient
	folderID         string
	server           string
	fileClient       *gofile.GofileClient
}

func NewScreenshot() *Screenshot {
	screenshotClient := screenshot.NewClient(os.Getenv("SCREENSHOT_API_KEY"))
	folderID := os.Getenv("GOFILE_FOLDER")
	server := os.Getenv("GOFILE_SERVER")
	fileClient := gofile.NewClient(os.Getenv("GOFILE_TOKEN"))
	return &Screenshot{
		screenshotClient: screenshotClient,
		folderID:         folderID,
		server:           server,
		fileClient:       fileClient,
	}
}

func (s *Screenshot) Create(ctx context.Context, url string) (*string, error) {
	var image_data []byte
	var err error
	if os.Getenv("SCREENSHOT_TYPE") == "Apiflash" {
		image_data, err = s.screenshotClient.Apiflash.GetScreenshot(ctx, url, map[string]string{"wait_until": "page_loaded", "format": "png"})
	} else {
		image_data, err = s.screenshotClient.Screenshotmachine.GetScreenshot(ctx, url, map[string]string{"dimension": "320x240", "delay": "2000", "format": "png"})
	}

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(image_data)
	fileName := fmt.Sprintf("%d.png", time.Now().UnixMilli())
	request := gofile.UploadFileRequest{
		File:     reader,
		FolderId: &s.folderID,
		Server:   &s.server,
		FileName: fileName,
	}

	resp, err := s.fileClient.Contents.UploadFile(ctx, request)
	if err != nil {
		return nil, err
	}

	fileUrl := fmt.Sprintf("https://%s.gofile.io/download/web/%s/%s", s.server, resp.Data.ID, fileName)
	return &fileUrl, nil
}
