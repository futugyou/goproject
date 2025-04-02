package types

 type Annotations struct {
	Audience  []Role `json:"audience"`
	Priority  *float32 `json:"priority"`
 }