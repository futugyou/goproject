package protocol

type PaginatedRequestParams struct {
	RequestParams `json:",inline"`
	Cursor        *string `json:"cursor"`
}
