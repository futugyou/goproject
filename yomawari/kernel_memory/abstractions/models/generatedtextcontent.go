package models

type GeneratedTextContent struct {
	Text       string
	TokenUsage *TokenUsage
}

func NewGeneratedTextContent(text string, tokenUsage *TokenUsage) *GeneratedTextContent {
	return &GeneratedTextContent{text, tokenUsage}
}

func (g *GeneratedTextContent) ToString() string {
	if g == nil {
		return ""
	}

	return g.Text
}
