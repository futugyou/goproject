package types

type ReadResourceRequestParams struct {
	RequestParams
	Uri *string `json:"uri"`
}
