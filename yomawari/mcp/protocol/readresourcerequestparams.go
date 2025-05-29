package protocol

type ReadResourceRequestParams struct {
	RequestParams `json:",inline"`
	Uri           string `json:"uri"`
}
