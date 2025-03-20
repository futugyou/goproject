package models

import "encoding/json"

type SearchResult struct {
	Query    string     `json:"query"`
	NoResult bool       `json:"noResult"`
	Results  []Citation `json:"results,omitempty"`
}

func (a *SearchResult) ToJson() string {
	if a == nil {
		return "{}"
	}

	if j, err := json.Marshal(a); err != nil {
		return "{}"
	} else {
		return string(j)
	}
}
