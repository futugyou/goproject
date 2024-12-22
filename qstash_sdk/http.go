package qstash

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpClient struct {
	http    *http.Client
	token   string
	baseurl string
}

func NewHttpClient(token string, baseUrl string) *httpClient {
	c := &httpClient{
		token:   token,
		baseurl: baseUrl,
		http:    &http.Client{},
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

func (c *httpClient) Delete(ctx context.Context, path string, request, response interface{}) error {
	return c.doRequest(ctx, path, "DELETE", nil, response)
}

func (c *httpClient) PostStream(ctx context.Context, path string, request interface{}) (*StreamResponse, error) {
	return c.doStreamRequest(ctx, path, "POST", request)
}

func (c *httpClient) doStreamRequest(ctx context.Context, path, method string, request interface{}) (*StreamResponse, error) {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.token))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if er := checkResponseStatusCode(resp); er != nil {
		return nil, er
	}

	streamResponse := &StreamResponse{
		Reader:    bufio.NewReader(resp.Body),
		Response:  resp,
		StreamEnd: false,
	}

	return streamResponse, nil
}

func (c *httpClient) doRequest(ctx context.Context, path, method string, request, response interface{}) error {
	path = c.createSubpath(path)
	var body io.Reader

	customeHeader := map[string]string{}
	if request != nil {
		if qstashReq, ok := request.(QstashRequest); ok {
			customeHeader = qstashReq.BuilderHeader()
			payload := qstashReq.GetPayload()
			if len(payload) > 0 {
				body = strings.NewReader(payload)
			}
		} else {
			payloadBytes, _ := json.Marshal(request)
			body = bytes.NewReader(payloadBytes)
		}
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.token))

	for key, value := range customeHeader {
		req.Header.Set(key, value)
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
