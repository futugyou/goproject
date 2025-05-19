package chat_completion

import (
	"strings"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions/ai_functional"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
)

const MessageTagName string = "message"
const RoleAttributeName string = "role"
const ImageTagName string = "image"
const TextTagName string = "text"

func ChatPromptParser(prompt string) (chatHistory *ChatHistory, ok bool) {
	messageTagStart := "<" + MessageTagName
	var nodes []ai_functional.PromptNode

	if !strings.Contains(prompt, messageTagStart) {
		return
	}

	nodes, ok = ai_functional.XmlPromptParser(prompt)
	if !ok {
		return
	}

	chatHistory, ok = chatPromptParser(nodes)
	return
}

func chatPromptParser(nodes []ai_functional.PromptNode) (*ChatHistory, bool) {
	chatHistory := &ChatHistory{}

	for _, node := range nodes {
		if isValidChatMessage(node) {
			chatHistory.Add(parseChatNode(node))
		}
	}

	return chatHistory, chatHistory.Count() > 0
}

func parseChatNode(node ai_functional.PromptNode) contents.ChatMessageContent {
	items := contents.ChatMessageContentItemCollection{}

	for _, childNode := range node.ChildNodes {
		if strings.Contains(childNode.TagName, ImageTagName) {
			if strings.HasPrefix(childNode.Content, "data:") {
				items.Add(contents.ImageContent{Data: []byte(childNode.Content)})
			} else {

				items.Add(contents.ImageContent{DataUri: childNode.Content})
			}
		} else if strings.Contains(childNode.TagName, TextTagName) {
			items.Add(contents.TextContent{Text: childNode.Content})
		}
	}

	if items.Count() == 1 {
		if con, ok := items.Get(0).(contents.TextContent); ok {
			node.Content = con.Text
			items.Clear()
		}
	}

	authorRole := contents.CreateAuthorRole(node.Attributes[RoleAttributeName])

	if items.Count() > 0 {
		return contents.ChatMessageContent{
			Role:  authorRole,
			Items: &items,
		}
	}
	return contents.ChatMessageContent{
		Role:    authorRole,
		Content: node.Content,
	}
}

func isValidChatMessage(node ai_functional.PromptNode) bool {
	_, ok := node.Attributes[RoleAttributeName]
	return ok && strings.EqualFold(node.TagName, MessageTagName)
}
