package contents

// TextContent represents text-based content.
type TextContent struct {
	AIContent `json:",inline"`
	Text      string `json:"text,omitempty"`
}
