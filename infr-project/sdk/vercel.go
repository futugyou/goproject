package sdk

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type VercelClient struct {
	baseurl string
	token   string
	http    *http.Client
}

func NewVercelClient(baseurl string, token string) *VercelClient {
	c := &VercelClient{
		http: &http.Client{},
	}
	c.baseurl = baseurl
	c.token = token
	return c
}

func (v *VercelClient) GetProjects() string {
	path := v.baseurl + "/v9/projects"
	var body io.Reader
	req, _ := http.NewRequest("GET", path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", v.token))
	resp, err := v.http.Do(req)

	if err != nil {
		log.Println(err.Error())
		return "[]"
	}

	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return "[]"
	}
	return string(all)
}

func (v *VercelClient) GetProjectEnv(project string) string {
	path := v.baseurl + "/v9/projects/" + project + "/env"
	var body io.Reader
	req, _ := http.NewRequest("GET", path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", v.token))
	resp, err := v.http.Do(req)

	if err != nil {
		log.Println(err.Error())
		return "[]"
	}

	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return "[]"
	}
	return string(all)
}
