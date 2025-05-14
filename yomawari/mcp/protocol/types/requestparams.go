package types

type RequestParams struct {
	Meta *RequestParamsMetadata `json:"_meta"`
}

type RequestParamsMetadata struct {
	ProgressToken *ProgressToken `json:"progressToken"`
}
