package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type IHttpClient interface {
	SetOrganization(organization string)
	SetBaseUrl(baseurl string)
	Get(path string, response interface{}) *OpenaiError
	Post(path string, request, response interface{}) *OpenaiError
	Delete(path string, response interface{}) *OpenaiError
	PostWithFile(path string, request, response interface{}) *OpenaiError
	PostStream(path string, request interface{}) (*StreamResponse, *OpenaiError)
	GetStream(path string) (*StreamResponse, *OpenaiError)
}

const baseUrl string = "https://api.openai.com/v1/"

type httpClient struct {
	http         *http.Client
	apikey       string
	organization string
	baseurl      string
}

func newHttpClient(apikey string) *httpClient {
	return &httpClient{
		apikey:       apikey,
		organization: "",
		baseurl:      baseUrl,
		http:         &http.Client{},
	}
}

func (c *httpClient) SetOrganization(organization string) {
	c.organization = organization
}

func (c *httpClient) SetBaseUrl(baseurl string) {
	c.baseurl = baseurl
}

func (c *httpClient) Post(path string, request, response interface{}) *OpenaiError {
	return c.doRequest(path, "POST", request, response)
}

func (c *httpClient) Get(path string, response interface{}) *OpenaiError {
	return c.doRequest(path, "GET", nil, response)
}

func (c *httpClient) Delete(path string, response interface{}) *OpenaiError {
	return c.doRequest(path, "DELETE", nil, response)
}

func (c *httpClient) doRequest(path, method string, request, response interface{}) *OpenaiError {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.apikey))

	if len(c.organization) > 0 {
		req.Header.Set("organization", c.organization)
	}

	return c.readHttpResponse(req, response)
}

func (c *httpClient) readHttpResponse(req *http.Request, response interface{}) *OpenaiError {
	resp, err := c.http.Do(req)

	if err != nil {
		return systemError(err.Error())
	}

	defer resp.Body.Close()

	if err := checkResponseStatusCode(resp); err != nil {
		return err
	}

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		return systemError(err.Error())
	}

	switch result := response.(type) {
	case *string:
		*result = string(all)
	default:
		if err = json.Unmarshal(all, response); err != nil {
			return systemError(err.Error())
		}
	}

	return nil
}

func checkResponseStatusCode(resp *http.Response) *OpenaiError {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		all, err := io.ReadAll(resp.Body)

		if err != nil {
			// raw error message
			return systemError(err.Error())
		}

		var apiError *OpenaiError
		if jsonError := json.Unmarshal(all, apiError); jsonError != nil {
			// raw error message
			return systemError(string(all))
		}

		if apiError == nil || len(apiError.ErrorMessage) == 0 {
			// raw error message
			return systemError(string(all))
		}

		return apiError
	}

	return nil
}

func (c *httpClient) PostWithFile(path string, request, response interface{}) *OpenaiError {
	path = c.createSubpath(path)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	reType := reflect.TypeOf(request)
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		return systemError("request must ptr")
	}

	v := reflect.ValueOf(request).Elem()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		if value.IsZero() {
			continue
		}

		structField := v.Type().Field(i)
		fieldName := structField.Name
		tag := structField.Tag.Get("json")
		if len(tag) > 0 {
			if strings.Contains(tag, ",") {
				fieldName = strings.Split(tag, ",")[0]
			} else {
				fieldName = tag
			}
		}

		switch v := value.Interface().(type) {
		case float32:
			writer.WriteField(fieldName, fmt.Sprintf("%f", v))
		case float64:
			writer.WriteField(fieldName, fmt.Sprintf("%f", v))
		case int:
			writer.WriteField(fieldName, strconv.Itoa(v))
		case int32:
			writer.WriteField(fieldName, strconv.FormatInt(int64(v), 10))
		case int64:
			writer.WriteField(fieldName, strconv.FormatInt(v, 10))
		case string:
			writer.WriteField(fieldName, v)
		case *os.File:
			if wimage, e := writer.CreateFormFile(fieldName, v.Name()); e != nil {
				return systemError(e.Error())
			} else {
				io.Copy(wimage, v)
			}

		default:
			writer.WriteField(fieldName, fmt.Sprintf("%v", v))
		}
	}

	writer.Close()

	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.apikey))

	if len(c.organization) > 0 {
		req.Header.Set("organization", c.organization)
	}

	return c.readHttpResponse(req, response)
}

func (c *httpClient) PostStream(path string, request interface{}) (*StreamResponse, *OpenaiError) {
	return c.doStreamRequest(path, "POST", request)
}

func (c *httpClient) GetStream(path string) (*StreamResponse, *OpenaiError) {
	return c.doStreamRequest(path, "GET", nil)
}

func (c *httpClient) doStreamRequest(path, method string, request interface{}) (*StreamResponse, *OpenaiError) {
	path = c.createSubpath(path)
	var body io.Reader

	if request != nil {
		payloadBytes, _ := json.Marshal(request)
		body = bytes.NewReader(payloadBytes)
	}

	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.apikey))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	if len(c.organization) > 0 {
		req.Header.Set("organization", c.organization)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, systemError(err.Error())
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

func (c *httpClient) createSubpath(path string) string {
	return c.baseurl + path
}
