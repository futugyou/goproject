package vercel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	http          *http.Client
	token         string
	baseurl       string
	customeHeader map[string]string
}

func newClient(token string, baseUrl string) *httpClient {
	c := &httpClient{
		token:   token,
		baseurl: baseUrl,
		http:    &http.Client{},
	}
	return c
}

func newClientWithHeader(baseUrl string, customeHeader map[string]string) *httpClient {
	c := &httpClient{
		baseurl:       baseUrl,
		http:          &http.Client{},
		customeHeader: customeHeader,
	}
	return c
}

func newClientWithHttp(baseUrl string, client *http.Client) *httpClient {
	if client == nil {
		client = &http.Client{}
	}
	c := &httpClient{
		baseurl: baseUrl,
		http:    client,
	}
	return c
}

func (c *httpClient) Post(ctx context.Context, path string, request, response interface{}) error {
	return c.DoRequest(ctx, path, "POST", request, response)
}

func (c *httpClient) Put(ctx context.Context, path string, request, response interface{}) error {
	return c.DoRequest(ctx, path, "PUT", request, response)
}

func (c *httpClient) Patch(ctx context.Context, path string, request, response interface{}) error {
	return c.DoRequest(ctx, path, "PATCH", request, response)
}

func (c *httpClient) Get(ctx context.Context, path string, response interface{}) error {
	return c.DoRequest(ctx, path, "GET", nil, response)
}

func (c *httpClient) Delete(ctx context.Context, path string, response interface{}) error {
	return c.DoRequest(ctx, path, "DELETE", nil, response)
}

func (c *httpClient) DeleteWithBody(ctx context.Context, path string, request, response interface{}) error {
	return c.DoRequest(ctx, path, "DELETE", request, response)
}

func (c *httpClient) DoRequest(ctx context.Context, path, method string, request, response interface{}) error {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if len(c.token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.token))
	}

	for key, value := range c.customeHeader {
		if len(value) > 0 {
			req.Header.Set(key, value)
		}
	}

	return c.readHttpResponse(ctx, req, response)
}

func (c *httpClient) readHttpResponse(ctx context.Context, req *http.Request, response interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)

	if err != nil {
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
		_, err := io.ReadAll(resp.Body)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *httpClient) createSubpath(path string) string {
	return c.baseurl + path
}
