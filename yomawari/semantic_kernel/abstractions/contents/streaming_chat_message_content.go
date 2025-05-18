package contents

import (
	"encoding/base64"
	"encoding/json"
)

type StreamingChatMessageContentItemCollection struct {
	Items []StreamingKernelContent `json:"items"`
}

func (c StreamingChatMessageContentItemCollection) MarshalJSON() ([]byte, error) {
	var rawItems []json.RawMessage
	for _, item := range c.Items {
		b, err := MarshalStreamingKernelContent(item)
		if err != nil {
			return nil, err
		}
		rawItems = append(rawItems, b)
	}
	return json.Marshal(map[string]any{
		"items": rawItems,
	})
}

func (c *StreamingChatMessageContentItemCollection) UnmarshalJSON(data []byte) error {
	var raw struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, item := range raw.Items {
		content, err := UnmarshalStreamingKernelContent(item)
		if err != nil {
			return err
		}
		c.Items = append(c.Items, content)
	}
	return nil
}

type StreamingChatMessageContent struct {
	ChoiceIndex  int                                       `json:"choiceIndex"`
	ModelId      string                                    `json:"modelId"`
	Metadata     map[string]any                            `json:"metadata"`
	InnerContent any                                       `json:"-"`
	Content      string                                    `json:"content"`
	Items        StreamingChatMessageContentItemCollection `json:"items"`
	AuthorName   string                                    `json:"authorName"`
	Role         AuthorRole                                `json:"role"`
	Encoding     *base64.Encoding                          `json:"-"`
}

func (StreamingChatMessageContent) Type() string {
	return "streamingChatMessage"
}

func (c StreamingChatMessageContent) ToString() string {
	return c.Content
}

func (c StreamingChatMessageContent) Hash() string {
	return c.Content
}

func (c StreamingChatMessageContent) ToByteArray() []byte {
	r, _ := base64.URLEncoding.DecodeString(c.ToString())
	return r
}

func (c StreamingChatMessageContent) GetContent() string {
	for _, item := range c.Items.Items {
		if textContent, ok := item.(StreamingTextContent); ok && item.Type() == "streaming-function-call-update" {
			return textContent.Text
		}
	}
	return ""
}

func (c *StreamingChatMessageContent) SetContent(content string) {
	for i, item := range c.Items.Items {
		if textContent, ok := item.(StreamingTextContent); ok && item.Type() == "streaming-function-call-update" {
			textContent.Text = content
			c.Items.Items[i] = textContent
			return
		}
	}

	var textContent StreamingTextContent = StreamingTextContent{
		ChoiceIndex:  c.ChoiceIndex,
		ModelId:      c.ModelId,
		Metadata:     c.Metadata,
		InnerContent: c.InnerContent,
		Text:         content,
		Encoding:     c.Encoding,
	}
	c.Items.Items = append(c.Items.Items, textContent)
}
