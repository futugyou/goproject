package types

type SubscribeRequestParams struct {
	RequestParams `json:",inline"`
	Uri           *string `json:"uri"`
}
