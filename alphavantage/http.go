package alphavantage

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
)

type httpClient struct {
	http *http.Client
}

func newHttpClient() *httpClient {
	return &httpClient{
		http: &http.Client{},
	}
}

func (c *httpClient) get(path string) (io.ReadCloser, error) {
	var body io.Reader
	req, _ := http.NewRequest("GET", path, body)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if err := checkResponseStatusCode(resp); err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func checkResponseStatusCode(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http request error")
	}

	return nil
}

func (c *httpClient) getCsv(path string) ([][]string, error) {
	readCloser, err := c.get(path)
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	reader := csv.NewReader(readCloser)
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	result, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}