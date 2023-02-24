package pkg

const BASEURL string = "https://api.openai.com/v1/"

type openaiclient struct {
	apikey       string
	organization string
	baseurl      string
}

func NewClient(apikey string, organization string, baseurl string) *openaiclient {
	if len(apikey) == 0 {
		panic("apikey can not be null")
	}
	if len(baseurl) == 0 {
		baseurl = BASEURL
	}
	return &openaiclient{
		apikey,
		organization,
		baseurl,
	}
}
