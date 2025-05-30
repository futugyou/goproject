package gofile

import "net/http"

type GofileClient struct {
	http     *http.Client
	key      string
	common   service
	Servers  *ServerService
	Contents *ContentService
	Accounts *AccountService
}

type service struct {
	client *GofileClient
}

// https://gofile.io/api
// A Premium account is required to get/search content, so the following APIs are skip.
func NewClient(token string) *GofileClient {
	c := &GofileClient{
		http: &http.Client{},
		key:  token,
	}
	c.initialize()
	return c
}

func (c *GofileClient) initialize() {
	c.common.client = c
	c.Servers = (*ServerService)(&c.common)
	c.Contents = (*ContentService)(&c.common)
	c.Accounts = (*AccountService)(&c.common)
}
