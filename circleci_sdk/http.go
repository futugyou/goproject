package circleci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type httpClient struct {
	http          *http.Client
	token         string
	baseurl       string
	customeHeader map[string]string
}

func NewHttpClient(token string, baseUrl string) *httpClient {
	c := &httpClient{
		token:   token,
		baseurl: baseUrl,
		http:    &http.Client{},
	}
	return c
}

func NewHttpClientWithHeader(baseUrl string, customeHeader map[string]string) *httpClient {
	c := &httpClient{
		baseurl:       baseUrl,
		http:          &http.Client{},
		customeHeader: customeHeader,
	}
	return c
}

func (c *httpClient) Post(ctx context.Context, path string, request, response interface{}) error {
	return c.doRequest(ctx, path, "POST", request, response)
}

func (c *httpClient) Put(ctx context.Context, path string, request, response interface{}) error {
	return c.doRequest(ctx, path, "PUT", request, response)
}

func (c *httpClient) Patch(ctx context.Context, path string, request, response interface{}) error {
	return c.doRequest(ctx, path, "PATCH", request, response)
}

func (c *httpClient) Get(ctx context.Context, path string, response interface{}) error {
	return c.doRequest(ctx, path, "GET", nil, response)
}

func (c *httpClient) Delete(ctx context.Context, path string, response interface{}) error {
	return c.doRequest(ctx, path, "DELETE", nil, response)
}

func (c *httpClient) doRequest(ctx context.Context, path, method string, request, response interface{}) error {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if len(c.customeHeader) > 0 {
		for key, value := range c.customeHeader {
			req.Header.Set(key, value)
		}
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.token))
	}

	return c.readHttpResponse(ctx, req, response)
}

func (c *httpClient) readHttpResponse(ctx context.Context, req *http.Request, response interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return e
			}
		}

		return err
	}

	defer resp.Body.Close()

	if err := checkResponseStatusCode(resp); err != nil {
		return err
	}

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	switch result := response.(type) {
	case *string:
		*result = string(all)
	default:
		if err = json.Unmarshal(all, response); err != nil {
			return err
		}
	}

	return nil
}

func checkResponseStatusCode(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("circleci service returned %d: %s", resp.StatusCode, string(data))
	}

	return nil
}

func (c *httpClient) createSubpath(path string) string {
	return c.baseurl + path
}
