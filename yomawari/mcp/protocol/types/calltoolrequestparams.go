package types

type CallToolRequestParams struct {
	RequestParams
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}
