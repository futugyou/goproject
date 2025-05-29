package protocol

type ToolAnnotations struct {
	Title           string `json:"title"`
	DestructiveHint *bool  `json:"destructiveHint"`
	IdempotentHint  *bool  `json:"idempotentHint"`
	OpenWorldHint   *bool  `json:"openWorldHint"`
	ReadOnlyHint    *bool  `json:"readOnlyHint"`
}
