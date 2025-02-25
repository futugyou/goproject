package screenshot

import (
	"context"
	"io"
	"net/http"
)

// https://www.screenshotmachine.com/website-screenshot-api.php
type ScreenshotmachineService service

const screenshotPath = "https://api.screenshotmachine.com"

func (s *ScreenshotmachineService) GetScreenshot(ctx context.Context, path string, header map[string]string) ([]byte, error) {
	request, _ := http.NewRequest("GET", screenshotPath, nil)
	data := request.URL.Query()
	data.Add("key", s.client.key)
	data.Add("url", path)

	if value, ok := header["dimension"]; ok {
		data.Add("dimension", value)
	}

	if value, ok := header["device"]; ok {
		data.Add("device", value)
	}

	if value, ok := header["format"]; ok {
		data.Add("format", value)
	}

	if value, ok := header["hash"]; ok {
		data.Add("hash", value)
	}

	if value, ok := header["cacheLimit"]; ok {
		data.Add("cacheLimit", value)
	}

	if value, ok := header["delay"]; ok {
		data.Add("delay", value)
	}

	if value, ok := header["zoom"]; ok {
		data.Add("zoom", value)
	}

	if value, ok := header["click"]; ok {
		data.Add("click", value)
	}

	if value, ok := header["hide"]; ok {
		data.Add("hide", value)
	}

	if value, ok := header["cookies"]; ok {
		data.Add("cookies", value)
	}

	if value, ok := header["accept-language"]; ok {
		data.Add("accept-language", value)
	}

	if value, ok := header["user-agent"]; ok {
		data.Add("user-agent", value)
	}

	if value, ok := header["selector"]; ok {
		data.Add("selector", value)
	}

	if value, ok := header["crop"]; ok {
		data.Add("crop", value)
	}

	request.URL.RawQuery = data.Encode()
	request = request.WithContext(ctx)
	resp, err := s.client.http.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
