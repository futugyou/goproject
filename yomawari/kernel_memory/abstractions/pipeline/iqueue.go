package pipeline

import "context"

type IQueue interface {
	ConnectToQueue(ctx context.Context, queueName string, options QueueOptions) (IQueue, error)
	Enqueue(ctx context.Context, message string) error
	OnDequeue(ctx context.Context, processMessageAction func(ctx context.Context, message string) ReturnType) error
}

type MessageHandler[T any] func(sender any, e T) error

type AsyncMessageHandler[T any] func(sender any, e T) <-chan error
