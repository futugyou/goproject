package types

type UnsubscribeRequestParams struct {
	RequestParams `json:",inline"`
	Uri           *string `json:"uri"`
}
