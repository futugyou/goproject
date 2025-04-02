package types

type CallToolResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError"`
}
