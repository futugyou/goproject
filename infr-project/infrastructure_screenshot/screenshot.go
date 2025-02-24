package infrastructure_screenshot

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
}

func (s *Screenshot) Create(ctx context.Context, url string) (*string, error) {
	sclient := screenshot.NewClient(os.Getenv("SCREENSHOT_API_KEY"))
	var image_data []byte
	var err error
	if os.Getenv("SCREENSHOT_TYPE") == "Apiflash" {
		image_data, err = sclient.Apiflash.GetScreenshot(ctx, url, map[string]string{"wait_until": "page_loaded", "format": "png"})
	} else {
		image_data, err = sclient.Screenshotmachine.GetScreenshot(ctx, url, map[string]string{"dimension": "320x240", "delay": "2000", "format": "png"})
	}

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(image_data)
	folderID := os.Getenv("GOFILE_FOLDER")
	server := os.Getenv("GOFILE_SERVER")
	fileName := fmt.Sprintf("%d.png", time.Now().UnixMilli())
	request := gofile.UploadFileRequest{
		File:     reader,
		FolderId: &folderID,
		Server:   &server,
		FileName: fileName,
	}

	fileClient := gofile.NewClient(os.Getenv("GOFILE_TOKEN"))
	resp, err := fileClient.Contents.UploadFile(ctx, request)
	if err != nil {
		return nil, err
	}

	fileUrl := fmt.Sprintf("https://%s.gofile.io/download/web/%s/%s", server, resp.Data.ID, fileName)
	return &fileUrl, nil
}
