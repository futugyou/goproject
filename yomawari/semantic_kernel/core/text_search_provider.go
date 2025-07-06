package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

var _ abstractions.AIContextProvider = (*TextSearchProvider)(nil)

const DefaultPluginSearchFunctionName string = "Search"
const DefaultPluginSearchFunctionDescription string = "Allows searching for additional information to help answer the user question."
const DefaultContextPrompt string = "## Additional Context\nConsider the following information from source documents when responding to the user:"
const DefaultIncludeCitationsPrompt string = "Include citations to the source document with document name and link if document name and link is available."

type TextSearchProvider struct {
	textSearch  abstractions.ITextSearch
	aIFunctions []functions.AIFunction
	Options     TextSearchProviderOptions
}

func NewTextSearchProvider(textSearch abstractions.ITextSearch, options TextSearchProviderOptions) *TextSearchProvider {
	provider := &TextSearchProvider{
		textSearch:  textSearch,
		aIFunctions: []functions.AIFunction{},
		Options:     options,
	}
	ops := &functions.AIFunctionFactoryOptions{}
	fun, err := functions.NewAIFunctionFactory().Create(provider.Search, ops)
	if err == nil {
		provider.aIFunctions = append(provider.aIFunctions, fun)
	}
	return provider
}

// ConversationCreated implements abstractions.AIContextProvider.
func (t *TextSearchProvider) ConversationCreated(ctx context.Context, conversationId string) error {
	return nil
}

// ConversationDeleting implements abstractions.AIContextProvider.
func (t *TextSearchProvider) ConversationDeleting(ctx context.Context, conversationId string) error {
	return nil
}

// MessageAdding implements abstractions.AIContextProvider.
func (t *TextSearchProvider) MessageAdding(ctx context.Context, conversationId string, newMessage chatcompletion.ChatMessage) error {
	return nil
}

// ModelInvoking implements abstractions.AIContextProvider.
func (t *TextSearchProvider) ModelInvoking(ctx context.Context, newMessages []chatcompletion.ChatMessage) (*abstractions.AIContext, error) {
	if t.Options.SearchTime != RagBehaviorBeforeAIInvoke {
		return &abstractions.AIContext{
			AIFunctions: t.aIFunctions,
		}, nil
	}

	inputs := []string{}
	for i := 0; i < len(newMessages); i++ {
		inputs = append(inputs, newMessages[i].Text())
	}
	input := strings.Join(inputs, "\n")
	ops := abstractions.TextSearchOptions{
		Top: t.Options.Top,
	}
	searchResults, err := t.textSearch.GetTextSearchResults(ctx, input, ops)
	if err != nil {
		return nil, err
	}
	formatted := t.formatResults(searchResults.Results)
	return &abstractions.AIContext{
		Instructions: formatted,
	}, nil
}

// Resuming implements abstractions.AIContextProvider.
func (t *TextSearchProvider) Resuming(ctx context.Context, conversationId string) error {
	return nil
}

// Suspending implements abstractions.AIContextProvider.
func (t *TextSearchProvider) Suspending(ctx context.Context, conversationId string) error {
	return nil
}

func (t *TextSearchProvider) Search(ctx context.Context, userQuestion string) (string, error) {
	searchResults, err := t.textSearch.GetTextSearchResults(ctx, userQuestion, abstractions.TextSearchOptions{
		Top: t.Options.Top,
	})
	if err != nil {
		return "", err
	}
	var formatted = t.formatResults(searchResults.Results)

	return formatted, nil
}

func (t *TextSearchProvider) formatResults(results []abstractions.TextSearchResult) string {
	if t.Options.ContextFormatter != nil {
		return t.Options.ContextFormatter(results)
	}

	if len(results) == 0 {
		return ""
	}

	var sb = &strings.Builder{}
	if len(t.Options.ContextPrompt) == 0 {
		sb.WriteString(DefaultContextPrompt)
		text.AppendLine(sb)
	} else {
		sb.WriteString(t.Options.ContextPrompt)
		text.AppendLine(sb)
	}

	for i := 0; i < len(results); i++ {
		var result = results[i]
		if len(result.Name) > 0 {
			sb.WriteString(fmt.Sprintf("SourceDocName: %s", result.Name))
			text.AppendLine(sb)
		}
		if len(result.Link) > 0 {
			sb.WriteString(fmt.Sprintf("SourceDocLink: %s", result.Link))
			text.AppendLine(sb)
		}
		sb.WriteString(fmt.Sprintf("Contents: %s", result.Value))
		text.AppendLine(sb)
		sb.WriteString("----")
		text.AppendLine(sb)
	}
	if len(t.Options.IncludeCitationsPrompt) == 0 {
		sb.WriteString(DefaultIncludeCitationsPrompt)
		text.AppendLine(sb)
	} else {
		sb.WriteString(t.Options.IncludeCitationsPrompt)
		text.AppendLine(sb)
	}

	text.AppendLine(sb)
	return sb.String()
}
