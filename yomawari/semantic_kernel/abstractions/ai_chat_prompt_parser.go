package abstractions

import (
	"strings"
)

const MessageTagName string = "message"
const RoleAttributeName string = "role"
const ImageTagName string = "image"
const TextTagName string = "text"

func ChatPromptParser(prompt string) (chatHistory *ChatHistory, ok bool) {
	messageTagStart := "<" + MessageTagName
	var nodes []PromptNode

	if !strings.Contains(prompt, messageTagStart) {
		return
	}

	nodes, ok = XmlPromptParser(prompt)
	if !ok {
		return
	}

	chatHistory, ok = chatPromptParser(nodes)
	return
}

func chatPromptParser(nodes []PromptNode) (*ChatHistory, bool) {
	chatHistory := &ChatHistory{}

	for _, node := range nodes {
		if isValidChatMessage(node) {
			chatHistory.Add(parseChatNode(node))
		}
	}

	return chatHistory, chatHistory.Count() > 0
}

func parseChatNode(node PromptNode) ChatMessageContent {
	items := ChatMessageContentItemCollection{}

	for _, childNode := range node.ChildNodes {
		if strings.Contains(childNode.TagName, ImageTagName) {
			if strings.HasPrefix(childNode.Content, "data:") {
				items.Add(ImageContent{Data: []byte(childNode.Content)})
			} else {

				items.Add(ImageContent{DataUri: childNode.Content})
			}
		} else if strings.Contains(childNode.TagName, TextTagName) {
			items.Add(TextContent{Text: childNode.Content})
		}
	}

	if items.Count() == 1 {
		if con, ok := items.Get(0).(TextContent); ok {
			node.Content = con.Text
			items.Clear()
		}
	}

	authorRole := CreateAuthorRole(node.Attributes[RoleAttributeName])

	if items.Count() > 0 {
		return ChatMessageContent{
			Role:  authorRole,
			Items: &items,
		}
	}
	return ChatMessageContent{
		Role:    authorRole,
		Content: node.Content,
	}
}

func isValidChatMessage(node PromptNode) bool {
	_, ok := node.Attributes[RoleAttributeName]
	return ok && strings.EqualFold(node.TagName, MessageTagName)
}
