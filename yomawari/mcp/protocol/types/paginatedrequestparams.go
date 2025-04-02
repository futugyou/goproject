package types

type PaginatedRequestParams struct {
	RequestParams
	Cursor *string `json:"cursor"`
}
