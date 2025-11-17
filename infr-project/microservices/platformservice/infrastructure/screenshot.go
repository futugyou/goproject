package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/futugyou/gofile"
	"github.com/futugyou/platformservice/options"
	"github.com/futugyou/screenshot"
)

type Screenshot struct {
	screenshotClient *screenshot.ScreenshotClient
	folderID         string
	server           string
	fileClient       *gofile.GofileClient
	opts             *options.Options
}

func NewScreenshot(opts *options.Options) *Screenshot {
	screenshotClient := screenshot.NewClient(opts.ScreenshotApiKey)
	fileClient := gofile.NewClient(opts.GofileToken)
	return &Screenshot{
		screenshotClient: screenshotClient,
		folderID:         opts.GofileFolder,
		server:           opts.GofileServer,
		fileClient:       fileClient,
		opts:             opts,
	}
}

func (s *Screenshot) Create(ctx context.Context, url string) (*string, error) {
	var image_data []byte
	var err error
	if s.opts.ScreenshotType == "Apiflash" {
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
