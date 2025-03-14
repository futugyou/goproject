package httputils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/futugyou/yomawari/common/errorutils"
)

const baseUrl string = "https://api.openai.com/v1/"

type HttpClient struct {
	http         *http.Client
	apikey       string
	organization string
	baseurl      string
}

func NewHttpClient(apikey string) *HttpClient {
	return &HttpClient{
		apikey:       apikey,
		organization: "",
		baseurl:      baseUrl,
		http:         &http.Client{},
	}
}

func (c *HttpClient) SetOrganization(organization string) {
	c.organization = organization
}

func (c *HttpClient) SetBaseUrl(baseurl string) {
	c.baseurl = baseurl
}

func (c *HttpClient) Post(ctx context.Context, path string, request, response interface{}) *errorutils.OpenaiError {
	return c.doRequest(ctx, path, "POST", request, response)
}

func (c *HttpClient) Get(ctx context.Context, path string, response interface{}) *errorutils.OpenaiError {
	return c.doRequest(ctx, path, "GET", nil, response)
}

func (c *HttpClient) Delete(ctx context.Context, path string, response interface{}) *errorutils.OpenaiError {
	return c.doRequest(ctx, path, "DELETE", nil, response)
}

func (c *HttpClient) doRequest(ctx context.Context, path, method string, request, response interface{}) *errorutils.OpenaiError {
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

	return c.readHttpResponse(ctx, req, response)
}

func (c *HttpClient) readHttpResponse(ctx context.Context, req *http.Request, response interface{}) *errorutils.OpenaiError {
	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)

	if err != nil {
		return errorutils.SystemError(err.Error())
	}

	defer resp.Body.Close()

	if err := checkResponseStatusCode(resp); err != nil {
		return err
	}

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		return errorutils.SystemError(err.Error())
	}

	switch result := response.(type) {
	case *string:
		*result = string(all)
	default:
		if err = json.Unmarshal(all, response); err != nil {
			return errorutils.SystemError(err.Error())
		}
	}

	return nil
}

func checkResponseStatusCode(resp *http.Response) *errorutils.OpenaiError {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		all, err := io.ReadAll(resp.Body)

		if err != nil {
			// raw error message
			return errorutils.SystemError(err.Error())
		}

		apiError := &errorutils.OpenaiError{}
		if jsonError := json.Unmarshal(all, apiError); jsonError != nil {
			// raw error message
			return errorutils.SystemError(string(all))
		}

		if len(apiError.ErrorMessage) == 0 {
			// raw error message
			return errorutils.SystemError(string(all))
		}

		return apiError
	}

	return nil
}

func (c *HttpClient) PostWithFile(ctx context.Context, path string, request, response interface{}) *errorutils.OpenaiError {
	path = c.createSubpath(path)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	reType := reflect.TypeOf(request)
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		return errorutils.SystemError("request must ptr")
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
				return errorutils.SystemError(e.Error())
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

	return c.readHttpResponse(ctx, req, response)
}

func (c *HttpClient) PostStream(ctx context.Context, path string, request interface{}) (*StreamResponse, *errorutils.OpenaiError) {
	return c.doStreamRequest(ctx, path, "POST", request)
}

func (c *HttpClient) GetStream(ctx context.Context, path string) (*StreamResponse, *errorutils.OpenaiError) {
	return c.doStreamRequest(ctx, path, "GET", nil)
}

func (c *HttpClient) doStreamRequest(ctx context.Context, path, method string, request interface{}) (*StreamResponse, *errorutils.OpenaiError) {
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

	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errorutils.SystemError(err.Error())
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

func (c *HttpClient) createSubpath(path string) string {
	return c.baseurl + path
}
