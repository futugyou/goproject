package vercel

type VercelClient struct {
	token string
	http  iHttpClient
}

const vercle_url string = "https://api.vercel.com"

func NewVercelClient(token string) *VercelClient {
	c := &VercelClient{
		http: newHttpClient(token, vercle_url),
	}
	c.token = token
	return c
}

func NewVercelClientWithHeader(customeHeader map[string]string) *VercelClient {
	c := &VercelClient{
		http: newHttpClientWithHeader(vercle_url, customeHeader),
	}
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
