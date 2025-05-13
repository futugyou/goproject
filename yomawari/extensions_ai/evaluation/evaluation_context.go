package evaluation

import "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"

type EvaluationContext interface {
	GetContents() []contents.IAIContent
	SetContents(conts []contents.IAIContent)
	GetName() string
	SetName(name string)
}

var _ EvaluationContext = (*BaseEvaluationContext)(nil)

type BaseEvaluationContext struct {
	name  string
	conts []contents.IAIContent
}

func NewBaseEvaluationContext(name string, conts []contents.IAIContent) *BaseEvaluationContext {
	return &BaseEvaluationContext{
		name:  name,
		conts: conts,
	}
}

// GetContents implements EvaluationContext.
func (b *BaseEvaluationContext) GetContents() []contents.IAIContent {
	return b.conts
}

// GetName implements EvaluationContext.
func (b *BaseEvaluationContext) GetName() string {
	return b.name
}

// SetContents implements EvaluationContext.
func (b *BaseEvaluationContext) SetContents(conts []contents.IAIContent) {
	b.conts = conts
}

// SetName implements EvaluationContext.
func (b *BaseEvaluationContext) SetName(name string) {
	b.name = name
}
