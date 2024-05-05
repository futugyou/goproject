package sdk

import "log"

type CircleciClient struct {
	http IHttpClient
}

const circleci_url string = "https://circleci.com/api/v2"

func NewCircleciClient(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: newHttpClientWithHeader(circleci_url, header),
	}
	return c
}

func (s *CircleciClient) Pipelines(org_slug string) string {
	path := "/pipeline?org-slug=" + org_slug
	result := "[]"
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}
