package dataformats

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type WebScraper struct {
	client *http.Client
}

func NewWebScraper(client *http.Client) *WebScraper {
	if client == nil {
		client = &http.Client{}
	}
	return &WebScraper{
		client: client,
	}
}

// GetContent implements dataformats.IWebScraper.
func (w *WebScraper) GetContent(ctx context.Context, url string) dataformats.WebScraperResult {
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return dataformats.WebScraperResult{Error: err, Success: false}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dataformats.WebScraperResult{Error: fmt.Errorf("HTTP status code %d", resp.StatusCode), Success: false}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return dataformats.WebScraperResult{Error: err, Success: false}
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return dataformats.WebScraperResult{Error: fmt.Errorf("no content type available"), Success: false}
	}

	mediaType := FixContentType(contentType, url)

	result := dataformats.WebScraperResult{
		Content:     data,
		ContentType: mediaType,
		Success:     false,
		Error:       err,
	}

	return result
}

func FixContentType(contentType string, path string) string {
	url, _ := url.Parse(path)
	if strings.Contains(contentType, pipeline.MimeTypes_PlainText) && strings.HasSuffix(url.Path, ".md") {
		return pipeline.MimeTypes_MarkDown
	}

	if strings.Contains(contentType, pipeline.MimeTypes_MarkDownOld1) || strings.Contains(contentType, pipeline.MimeTypes_MarkDownOld2) {
		return pipeline.MimeTypes_MarkDown
	}

	if strings.Contains(contentType, pipeline.MimeTypes_XML2) {
		return pipeline.MimeTypes_XML
	}

	mediaType, _, err := parseContentType(contentType)
	if err == nil {
		return mediaType
	}

	return contentType
}

func parseContentType(contentType string) (string, string, error) {
	parts := strings.Split(contentType, ";")
	if len(parts) > 0 {
		return parts[0], parts[1], nil
	}
	return contentType, "", nil
}

var _ dataformats.IWebScraper = (*WebScraper)(nil)
