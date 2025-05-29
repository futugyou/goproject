package protocol

type Completion struct {
	Values  []string `json:"values"`
	Total   *int     `json:"total"`
	HasMore *bool    `json:"hasMore"`
}
