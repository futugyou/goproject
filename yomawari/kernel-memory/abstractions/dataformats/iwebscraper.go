package dataformats

import "context"

type IWebScraper interface {
	GetContent(ctx context.Context, url string) WebScraperResult
}

type WebScraperResult struct {
	Content     []byte
	ContentType string
	Success     bool
	Error       error
}
