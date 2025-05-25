package abstractions

import (
	"html"
	"strings"

	"github.com/beevik/etree"
)

func XmlPromptParser(prompt string) ([]PromptNode, bool) {
	if prompt == "" || !strings.Contains(prompt, "<") || !(strings.Contains(prompt, "</") || strings.Contains(prompt, "/>")) {
		return nil, false
	}

	doc := etree.NewDocument()

	if err := doc.ReadFromString("<root>" + prompt + "</root>"); err != nil {
		return nil, false
	}

	var result []PromptNode
	for _, elem := range doc.Root().ChildElements() {
		if node := getPromptNode(elem); node != nil {
			result = append(result, *node)
		}
	}

	return result, len(result) > 0
}

func getPromptNode(elem *etree.Element) *PromptNode {
	node := &PromptNode{
		TagName:    elem.Tag,
		Attributes: map[string]string{},
	}

	for _, attr := range elem.Attr {
		node.Attributes[attr.Key] = attr.Value
	}

	trimmedText := strings.TrimSpace(elem.Text())
	if trimmedText != "" {
		node.Content = html.UnescapeString(trimmedText)
	}

	for _, child := range elem.ChildElements() {
		if childNode := getPromptNode(child); childNode != nil {
			node.ChildNodes = append(node.ChildNodes, *childNode)
		}
	}

	return node
}
