package evaluation

import "github.com/futugyou/yomawari/extensions-ai/abstractions/contents"

type EvaluationContext interface {
	GetContents() []contents.IAIContent
}
