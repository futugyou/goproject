package types

type CallToolRequestParams struct {
	RequestParams `json:",inline"`
	Name          string                 `json:"name"`
	Arguments     map[string]interface{} `json:"arguments"`
}
