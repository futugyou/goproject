package chat_completion

import "github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"

type ChatHistory struct {
	data []contents.ChatMessageContent
}

func (ge *ChatHistory) Count() int {
	return len(ge.data)
}

func (ge *ChatHistory) Add(item contents.ChatMessageContent) {
	ge.data = append(ge.data, item)
}

func (ge *ChatHistory) AddRange(items []contents.ChatMessageContent) {
	ge.data = append(ge.data, items...)
}

func (ge *ChatHistory) Clear() {
	ge.data = nil
}

func (ge *ChatHistory) Contains(item contents.ChatMessageContent) bool {
	for _, d := range ge.data {
		if d.Content == item.Content {
			return true
		}
	}
	return false
}

func (ge *ChatHistory) Get(index int) contents.ChatMessageContent {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	return ge.data[index]
}

func (ge *ChatHistory) Set(index int, item contents.ChatMessageContent) {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	ge.data[index] = item
}

func (ge *ChatHistory) Remove(item contents.ChatMessageContent) bool {
	for i, d := range ge.data {
		if d.Content == item.Content {
			ge.data = append(ge.data[:i], ge.data[i+1:]...)
			return true
		}
	}
	return false
}

func (ge *ChatHistory) RemoveAt(index int) {
	if index < 0 || index >= len(ge.data) {
		panic("index out of bounds")
	}
	ge.data = append(ge.data[:index], ge.data[index+1:]...)
}

func (ge *ChatHistory) IndexOf(item contents.ChatMessageContent) int {
	for i, d := range ge.data {
		if d.Content == item.Content {
			return i
		}
	}
	return -1
}

func (ge *ChatHistory) GetEnumerator() <-chan contents.ChatMessageContent {
	ch := make(chan contents.ChatMessageContent)
	go func() {
		defer close(ch)
		for _, embedding := range ge.data {
			ch <- embedding
		}
	}()
	return ch
}
