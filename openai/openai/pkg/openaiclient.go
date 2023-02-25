package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
		apikey,
		"",
		baseUrl,
	}
}

func (c *openaiClient) SetOrganization(organization string) {
	c.organization = organization
}

func (c *openaiClient) SetBaseUrl(baseurl string) {
	c.baseurl = baseurl
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
