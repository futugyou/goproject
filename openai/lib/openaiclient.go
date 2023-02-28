package lib

import (
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

const baseUrl string = "https://api.openai.com/v1/"

type openaiClient struct {
	apikey       string
	organization string
	baseurl      string
}

func NewClient(apikey string) *openaiClient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}
	return &openaiClient{
		apikey:       apikey,
		organization: "",
		baseurl:      baseUrl,
	}
}

func (c *openaiClient) SetOrganization(organization string) {
	c.organization = organization
}

func (c *openaiClient) SetBaseUrl(baseurl string) {
	c.baseurl = baseurl
}

func (c *openaiClient) Post(path string, request, response interface{}) {
	c.doRequest(path, "POST", request, response)
}

func (c *openaiClient) Get(path string, response interface{}) {
	c.doRequest(path, "GET", nil, response)
}

func (c *openaiClient) Delete(path string, response interface{}) {
	c.doRequest(path, "DELETE", nil, response)
}

func (c *openaiClient) doRequest(path, method string, request, response interface{}) {
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

	resp, err := http.DefaultClient.Do(req)

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

	json.Unmarshal(all, response)
}

func (c *openaiClient) PostWithFile(path string, request, response interface{}) {
	path = c.baseurl + path
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	reType := reflect.TypeOf(request)
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		fmt.Println("request must ptr")
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
			wimage, _ := writer.CreateFormFile(fieldName, v.Name())
			io.Copy(wimage, v)
		}
	}

	writer.Close()

	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", c.apikey))

	if len(c.organization) > 0 {
		req.Header.Set("organization", c.organization)
	}

	resp, err := http.DefaultClient.Do(req)

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

	json.Unmarshal(all, response)
}
