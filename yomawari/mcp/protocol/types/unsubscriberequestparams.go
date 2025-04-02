package types

type UnsubscribeRequestParams struct {
	RequestParams
	Uri *string `json:"uri"`
}
