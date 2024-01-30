package main

import (
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	http *http.Client
}

func NewHttpClient() *httpClient {
	return &httpClient{
		http: &http.Client{},
	}
}

func (c *httpClient) Get(path string) (interface{}, error) {
	var body io.Reader
	req, _ := http.NewRequest("GET", path, body)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := checkResponseStatusCode(resp); err != nil {
		return nil, err
	}

	response, err := ReadTimeSeries(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func checkResponseStatusCode(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http request error")
	}

	return nil
}
