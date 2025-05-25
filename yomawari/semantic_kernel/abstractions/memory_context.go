package abstractions

import (
	"context"
	"strings"
	"sync"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

type AIContext struct {
	Instructions string
	AIFunctions  []functions.AIFunction
}

type AIContextProvider interface {
	ConversationCreated(ctx context.Context, conversationId string) error
	MessageAdding(ctx context.Context, conversationId string, newMessage chatcompletion.ChatMessage) error
	ConversationDeleting(ctx context.Context, conversationId string) error
	ModelInvoking(ctx context.Context, newMessages []chatcompletion.ChatMessage) (*AIContext, error)
	Suspending(ctx context.Context, conversationId string) error
	Resuming(ctx context.Context, conversationId string) error
}

type AggregateAIContextProvider struct {
	Providers []AIContextProvider
}

func NewAggregateAIContextProvider(providers ...AIContextProvider) *AggregateAIContextProvider {
	return &AggregateAIContextProvider{Providers: providers}
}

func (a *AggregateAIContextProvider) Add(aiContextProvider AIContextProvider) {
	a.Providers = append(a.Providers, aiContextProvider)
}

func (a *AggregateAIContextProvider) AddFromServiceProvider(serviceProvider core.IServiceProvider) {
	p, _ := core.GetService[AIContextProvider](serviceProvider)
	if p != nil {
		a.Providers = append(a.Providers, p)
	}
}

func (a *AggregateAIContextProvider) ConversationCreated(ctx context.Context, conversationId string) error {
	return runParallel(a.Providers, func(p AIContextProvider) error {
		return p.ConversationCreated(ctx, conversationId)
	})
}

// ConversationDeleting implements AIContextProvider.
func (a *AggregateAIContextProvider) ConversationDeleting(ctx context.Context, conversationId string) error {
	return runParallel(a.Providers, func(p AIContextProvider) error {
		return p.ConversationDeleting(ctx, conversationId)
	})
}

// MessageAdding implements AIContextProvider.
func (a *AggregateAIContextProvider) MessageAdding(ctx context.Context, conversationId string, newMessage chatcompletion.ChatMessage) error {
	return runParallel(a.Providers, func(p AIContextProvider) error {
		return p.MessageAdding(ctx, conversationId, newMessage)
	})
}

// Resuming implements AIContextProvider.
func (a *AggregateAIContextProvider) Resuming(ctx context.Context, conversationId string) error {
	return runParallel(a.Providers, func(p AIContextProvider) error {
		return p.Resuming(ctx, conversationId)
	})
}

// Suspending implements AIContextProvider.
func (a *AggregateAIContextProvider) Suspending(ctx context.Context, conversationId string) error {
	return runParallel(a.Providers, func(p AIContextProvider) error {
		return p.Suspending(ctx, conversationId)
	})
}

func runParallel[T any](items []T, fn func(T) error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(items))

	for _, item := range items {
		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			if err := fn(item); err != nil {
				errChan <- err
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

// ModelInvoking implements AIContextProvider.
func (a *AggregateAIContextProvider) ModelInvoking(ctx context.Context, newMessages []chatcompletion.ChatMessage) (*AIContext, error) {
	var wg sync.WaitGroup
	var err error
	var aiContext *AIContext
	subContexts := []AIContext{}
	for _, provider := range a.Providers {
		wg.Add(1)
		go func(provider AIContextProvider) {
			defer wg.Done()
			aiContext, err = provider.ModelInvoking(ctx, newMessages)
			if err != nil {
				return
			}
			if aiContext != nil {
				subContexts = append(subContexts, *aiContext)
			}
		}(provider)
	}
	wg.Wait()

	if err != nil {
		return nil, err
	}

	aiContext = &AIContext{}
	instructions := []string{}
	for _, v := range subContexts {
		if len(v.AIFunctions) > 0 {
			aiContext.AIFunctions = append(aiContext.AIFunctions, v.AIFunctions...)
		}
		if len(v.Instructions) > 0 {
			instructions = append(instructions, v.Instructions)
		}
	}
	aiContext.Instructions = strings.Join(instructions, "\n")

	return aiContext, nil
}

func (a *AggregateAIContextProvider) MessageContentAdding(ctx context.Context, conversationId string, newMessage ChatMessageContent) error {
	return a.MessageAdding(ctx, conversationId, newMessage.ToChatMessage())
}

func (a *AggregateAIContextProvider) ContentModelInvoking(ctx context.Context, newMessages []ChatMessageContent) (*AIContext, error) {
	messages := []chatcompletion.ChatMessage{}
	for _, v := range newMessages {
		messages = append(messages, v.ToChatMessage())
	}
	return a.ModelInvoking(ctx, messages)
}
