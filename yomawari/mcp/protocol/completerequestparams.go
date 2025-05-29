package protocol

type CompleteRequestParams struct {
	Ref      Reference `json:"ref"`
	Argument Argument  `json:"argument"`
}
