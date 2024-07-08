package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IHttpClient interface {
	Get(path string, response interface{}) error
	Post(path string, request, response interface{}) error
	Put(path string, request, response interface{}) error
	Patch(path string, request, response interface{}) error
	Delete(path string, response interface{}) error
	DeleteWithBody(path string, request, response interface{}) error
}

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

func (c *httpClient) Post(path string, request, response interface{}) error {
	return c.doRequest(path, "POST", request, response)
}

func (c *httpClient) Put(path string, request, response interface{}) error {
	return c.doRequest(path, "PUT", request, response)
}

func (c *httpClient) Patch(path string, request, response interface{}) error {
	return c.doRequest(path, "PATCH", request, response)
}

func (c *httpClient) Get(path string, response interface{}) error {
	return c.doRequest(path, "GET", nil, response)
}

func (c *httpClient) Delete(path string, response interface{}) error {
	return c.doRequest(path, "DELETE", nil, response)
}

func (c *httpClient) DeleteWithBody(path string, request, response interface{}) error {
	return c.doRequest(path, "DELETE", request, response)
}

func (c *httpClient) doRequest(path, method string, request, response interface{}) error {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	if c.customeHeader != nil && len(c.customeHeader) > 0 {
		for key, value := range c.customeHeader {
			req.Header.Set(key, value)
		}
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.token))
	}

	return c.readHttpResponse(req, response)
}

func (c *httpClient) readHttpResponse(req *http.Request, response interface{}) error {
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
