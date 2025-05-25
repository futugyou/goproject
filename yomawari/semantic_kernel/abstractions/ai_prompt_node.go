package abstractions

type PromptNode struct {
	TagName    string
	Content    string
	Attributes map[string]string
	ChildNodes []PromptNode
}

func NewPromptNode(tagName string) *PromptNode {
	return &PromptNode{
		TagName:    tagName,
		Attributes: make(map[string]string),
		ChildNodes: make([]PromptNode, 0),
	}
}
