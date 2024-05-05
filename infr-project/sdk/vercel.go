package sdk

import (
	"log"
)

type VercelClient struct {
	baseurl string
	token   string
	http    IHttpClient
}

func NewVercelClient(baseurl string, token string) *VercelClient {
	c := &VercelClient{
		http: newHttpClient(token, baseurl),
	}
	c.baseurl = baseurl
	c.token = token
	return c
}

func (v *VercelClient) GetProjects() string {
	path := v.baseurl + "/v9/projects"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (v *VercelClient) GetProjectEnv(project string) string {
	path := v.baseurl + "/v9/projects/" + project + "/env"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}
