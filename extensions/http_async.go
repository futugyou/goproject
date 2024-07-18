package extensions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IAsyncHttpClient interface {
	GetAsync(path string) (<-chan *HttpResponse, <-chan error)
	PostAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error)
	PutAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error)
	PatchAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error)
	DeleteAsync(path string) (<-chan *HttpResponse, <-chan error)
	DeleteWithBodyAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error)
}

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

// Implementing the IAsyncHttpClient interface

func (c *httpClient) GetAsync(path string) (<-chan *HttpResponse, <-chan error) {
	respChan := make(chan *HttpResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		req, _ := http.NewRequest("GET", c.createSubpath(path), nil)
		req.Header.Set("Content-Type", "application/json")
		if c.customeHeader != nil && len(c.customeHeader) > 0 {
			for key, value := range c.customeHeader {
				req.Header.Set(key, value)
			}
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		}

		resp, err := c.http.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if err := checkResponseStatusCode(resp); err != nil {
			errChan <- err
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
			return
		}

		respChan <- &HttpResponse{StatusCode: resp.StatusCode, Body: body}
	}()

	return respChan, errChan
}

func (c *httpClient) PostAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error) {
	return c.doRequestAsync(path, "POST", request)
}

func (c *httpClient) PutAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error) {
	return c.doRequestAsync(path, "PUT", request)
}

func (c *httpClient) PatchAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error) {
	return c.doRequestAsync(path, "PATCH", request)
}

func (c *httpClient) DeleteAsync(path string) (<-chan *HttpResponse, <-chan error) {
	return c.doRequestAsync(path, "DELETE", nil)
}

func (c *httpClient) DeleteWithBodyAsync(path string, request interface{}) (<-chan *HttpResponse, <-chan error) {
	return c.doRequestAsync(path, "DELETE", request)
}

func (c *httpClient) doRequestAsync(path, method string, request interface{}) (<-chan *HttpResponse, <-chan error) {
	respChan := make(chan *HttpResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		var body io.Reader
		if request != nil {
			payloadBytes, _ := json.Marshal(request)
			body = bytes.NewReader(payloadBytes)
		}

		req, _ := http.NewRequest(method, c.createSubpath(path), body)
		req.Header.Set("Content-Type", "application/json")
		if c.customeHeader != nil && len(c.customeHeader) > 0 {
			for key, value := range c.customeHeader {
				req.Header.Set(key, value)
			}
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		}

		resp, err := c.http.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if err := checkResponseStatusCode(resp); err != nil {
			errChan <- err
			return
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
			return
		}

		respChan <- &HttpResponse{StatusCode: resp.StatusCode, Body: bodyBytes}
	}()

	return respChan, errChan
}
