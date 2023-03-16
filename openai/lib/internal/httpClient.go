package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	Get(path string, response interface{}) error
	Post(path string, request, response interface{}) error
	Delete(path string, response interface{}) error
	PostWithFile(path string, request, response interface{}) error
	PostStream(path string, request interface{}) error
	GetStream(path string) error
	ReadStream(response interface{})
	Close()
	CanReadStream() bool
}

const baseUrl string = "https://api.openai.com/v1/"
const endTag string = "[DONE]"

var headerData []byte = []byte("data: ")

type httpClient struct {
	http           *http.Client
	apikey         string
	organization   string
	baseurl        string
	streamResponse *http.Response
	streamReader   *bufio.Reader
	StreamEnd      bool
}

func NewHttpClient(apikey string) *httpClient {
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

func (c *httpClient) Post(path string, request, response interface{}) error {
	return c.doRequest(path, "POST", request, response)
}

func (c *httpClient) Get(path string, response interface{}) error {
	return c.doRequest(path, "GET", nil, response)
}

func (c *httpClient) Delete(path string, response interface{}) error {
	return c.doRequest(path, "DELETE", nil, response)
}

func (c *httpClient) doRequest(path, method string, request, response interface{}) (err error) {
	path = c.baseurl + path
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

	resp, err := c.http.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	switch result := response.(type) {
	case *string:
		*result = string(all)
	default:
		err = json.Unmarshal(all, response)
	}
	return
}

func (c *httpClient) PostWithFile(path string, request, response interface{}) (err error) {
	path = c.baseurl + path
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	reType := reflect.TypeOf(request)
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("request must ptr")
		return
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
				err = e
				return
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

	resp, err := c.http.Do(req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	switch result := response.(type) {
	case *string:
		*result = string(all)
	default:
		err = json.Unmarshal(all, response)
	}

	return
}

func (c *httpClient) PostStream(path string, request interface{}) error {
	return c.doStreamRequest(path, "POST", request)
}

func (c *httpClient) GetStream(path string) error {
	return c.doStreamRequest(path, "GET", nil)
}

func (c *httpClient) doStreamRequest(path, method string, request interface{}) (err error) {
	path = c.baseurl + path
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
		log.Println(err.Error())
		return
	}

	c.streamReader = bufio.NewReader(resp.Body)
	c.streamResponse = resp
	return
}

func (c *httpClient) ReadStream(response interface{}) {
	reader := c.streamReader

	if reader == nil {
		c.StreamEnd = true
		return
	}

	line, err := reader.ReadBytes('\n')
	responseStr := ""

	for {
		if err != nil {
			c.StreamEnd = true
			return
		}

		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, headerData) {
			line = bytes.TrimPrefix(line, headerData)
			responseStr = string(line)
			break
		} else {
			line, err = reader.ReadBytes('\n')
		}
	}

	if responseStr == endTag {
		c.StreamEnd = true
		return
	}

	json.Unmarshal(line, response)
}

func (c *httpClient) Close() {
	if c.streamResponse != nil {
		c.streamResponse.Body.Close()
	}
}

func (c *httpClient) CanReadStream() bool {
	return !c.StreamEnd
}
