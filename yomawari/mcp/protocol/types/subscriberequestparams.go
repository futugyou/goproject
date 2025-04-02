package types

type SubscribeRequestParams struct {
	RequestParams
	Uri *string `json:"uri"`
}
