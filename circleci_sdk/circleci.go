package circleci

type CircleciClient struct {
	http IHttpClient
}

const circleci_url string = "https://circleci.com/api/v2"

func NewCircleciClient(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: NewHttpClientWithHeader(circleci_url, header),
	}
	return c
}

type BaseResponse struct {
	Message *string `json:"message,omitempty"`
}
