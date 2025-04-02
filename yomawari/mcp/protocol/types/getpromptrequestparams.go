package types

type GetPromptRequestParams struct {
	RequestParams
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}
