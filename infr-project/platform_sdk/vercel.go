package sdk

import (
	"log"
)

type VercelClient struct {
	token string
	http  IHttpClient
}

const vercle_url string = "https://api.vercel.com"

func NewVercelClient(token string) *VercelClient {
	c := &VercelClient{
		http: newHttpClient(token, vercle_url),
	}
	c.token = token
	return c
}

func (v *VercelClient) GetProjects() string {
	path := "/v9/projects"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (v *VercelClient) GetProjectEnv(project string) string {
	path := "/v9/projects/" + project + "/env"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}
