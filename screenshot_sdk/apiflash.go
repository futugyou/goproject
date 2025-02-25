package screenshot

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// https://apiflash.com/documentation
type ApiflashService service

const flashPath = "https://api.apiflash.com/v1/urltoimage"

func (s *ApiflashService) GetScreenshot(ctx context.Context, path string, header map[string]string) ([]byte, error) {
	data := url.Values{}
	data.Add("access_key", s.client.key)
	data.Add("url", path)

	if value, ok := header["format"]; ok {
		data.Add("format", value)
	}

	if value, ok := header["width"]; ok {
		data.Add("width", value)
	}

	if value, ok := header["height"]; ok {
		data.Add("height", value)
	}

	if value, ok := header["fresh"]; ok {
		data.Add("fresh", value)
	}

	if value, ok := header["full_page"]; ok {
		data.Add("full_page", value)
	}

	if value, ok := header["quality"]; ok {
		data.Add("quality", value)
	}

	if value, ok := header["delay"]; ok {
		data.Add("delay", value)
	}

	if value, ok := header["scroll_page"]; ok {
		data.Add("scroll_page", value)
	}

	if value, ok := header["ttl"]; ok {
		data.Add("ttl", value)
	}

	if value, ok := header["response_type"]; ok {
		data.Add("response_type", value)
	}

	if value, ok := header["thumbnail_width"]; ok {
		data.Add("thumbnail_width", value)
	}

	if value, ok := header["crop"]; ok {
		data.Add("crop", value)
	}

	if value, ok := header["no_cookie_banners"]; ok {
		data.Add("no_cookie_banners", value)
	}

	if value, ok := header["no_ads"]; ok {
		data.Add("no_ads", value)
	}

	if value, ok := header["no_tracking"]; ok {
		data.Add("no_tracking", value)
	}

	if value, ok := header["scale_factor"]; ok {
		data.Add("scale_factor", value)
	}

	if value, ok := header["element"]; ok {
		data.Add("element", value)
	}

	if value, ok := header["element_overlap"]; ok {
		data.Add("element_overlap", value)
	}

	if value, ok := header["user_agent"]; ok {
		data.Add("user_agent", value)
	}

	if value, ok := header["extract_html"]; ok {
		data.Add("extract_html", value)
	}

	if value, ok := header["extract_text"]; ok {
		data.Add("extract_text", value)
	}

	if value, ok := header["transparent"]; ok {
		data.Add("transparent", value)
	}

	if value, ok := header["wait_for"]; ok {
		data.Add("wait_for", value)
	}

	if value, ok := header["wait_until"]; ok {
		data.Add("wait_until", value)
	}

	if value, ok := header["fail_on_status"]; ok {
		data.Add("fail_on_status", value)
	}

	if value, ok := header["accept_language"]; ok {
		data.Add("accept_language", value)
	}

	if value, ok := header["css"]; ok {
		data.Add("css", value)
	}

	if value, ok := header["cookies"]; ok {
		data.Add("cookies", value)
	}

	if value, ok := header["proxy"]; ok {
		data.Add("proxy", value)
	}

	if value, ok := header["latitude"]; ok {
		data.Add("latitude", value)
	}

	if value, ok := header["longitude"]; ok {
		data.Add("longitude", value)
	}

	if value, ok := header["accuracy"]; ok {
		data.Add("accuracy", value)
	}

	if value, ok := header["js"]; ok {
		data.Add("js", value)
	}

	if value, ok := header["headers"]; ok {
		data.Add("headers", value)
	}

	if value, ok := header["time_zone"]; ok {
		data.Add("time_zone", value)
	}

	if value, ok := header["ip_location"]; ok {
		data.Add("ip_location", value)
	}

	if value, ok := header["s3_access_key_id"]; ok {
		data.Add("s3_access_key_id", value)
	}

	if value, ok := header["s3_secret_key"]; ok {
		data.Add("s3_secret_key", value)
	}

	if value, ok := header["s3_bucket"]; ok {
		data.Add("s3_bucket", value)
	}

	if value, ok := header["s3_key"]; ok {
		data.Add("s3_key", value)
	}

	if value, ok := header["s3_endpoint"]; ok {
		data.Add("s3_endpoint", value)
	}

	if value, ok := header["s3_region"]; ok {
		data.Add("s3_region", value)
	}

	req, err := http.NewRequest("POST", flashPath, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)
	resp, err := s.client.http.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
