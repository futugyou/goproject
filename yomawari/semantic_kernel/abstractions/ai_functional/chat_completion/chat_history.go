package chat_completion

import (
	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
)

type ChatHistory struct {
	core.List[contents.ChatMessageContent]
}
