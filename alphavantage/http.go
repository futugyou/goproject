package alphavantage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/jszwec/csvutil"
)

type httpClient struct {
	http *http.Client
}

func newHttpClient() *httpClient {
	return &httpClient{
		http: &http.Client{},
	}
}

var errorFeilds = []string{
	"Error Message",
	"Information",
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
	fmt.Println(path)
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

	// {
	// 		"Information": "Thank you for using Alpha Vantage! Our standard API rate limit is 25 requests per day. Please subscribe to any of the premium plans at https://www.alphavantage.co/premium/ to instantly remove all daily rate limits."
	//		"Error Message": "Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for XXXX."
	// }
	if len(result) == 2 && len(result[0]) == 1 {
		for _, field := range errorFeilds {
			if message, ok := strings.CutPrefix(result[0][0], fmt.Sprintf("%s\": \"", field)); ok {
				return nil, fmt.Errorf(message)
			}
		}
	}

	return result, nil
}

func (c *httpClient) getCsvByUtil(path string, response interface{}) error {
	fmt.Println(path)
	readCloser, err := c.get(path)
	if err != nil {
		return err
	}

	defer readCloser.Close()

	reader := csv.NewReader(readCloser)
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	dec, err := csvutil.NewDecoder(reader)

	if err != nil {
		return err
	}

	h := dec.Header()
	if len(h) == 1 && h[0] == "{" {
		return fmt.Errorf("invalid API call or API rate limit is 25 requests per day")
	}

	timeFunc := csvutil.UnmarshalFunc(unmarshalTime)
	floatFunc := csvutil.UnmarshalFunc(unmarshalFloat)
	dec.WithUnmarshalers(
		csvutil.NewUnmarshalers(
			timeFunc,
			floatFunc,
		),
	)
	return dec.Decode(response)
}

func (c *httpClient) getJson(path string, response interface{}) error {
	fmt.Println(path)
	readCloser, err := c.get(path)
	if err != nil {
		return err
	}

	defer readCloser.Close()

	all, err := io.ReadAll(readCloser)
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

		// 1. respose is &sometype{}
		// 2. use reflect.Indirect or Elem
		ps := reflect.Indirect(reflect.ValueOf(response))
		// ps := reflect.ValueOf(response).Elem()
		for _, field := range errorFeilds {
			msg := ps.FieldByName(field).String()
			if len(msg) > 0 && msg != "<invalid Value>" {
				return fmt.Errorf(msg)
			}
		}
	}

	return nil
}
