package vercel

type VercelClient struct {
	token string
	http  IHttpClient
}

const vercle_url string = "https://api.vercel.com"

func NewVercelClient(token string) *VercelClient {
	c := &VercelClient{
		http: NewHttpClient(token, vercle_url),
	}
	c.token = token
	return c
}

type VercelError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Count int    `json:"count"`
	Next  string `json:"next"`
}
