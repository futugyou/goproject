package pipeline

import "context"

type QueueClientFactory struct {
	queueBuilder func(ctx context.Context) (IQueue, error)
}

func NewQueueClientFactory(queueBuilder func(ctx context.Context) (IQueue, error)) *QueueClientFactory {
	return &QueueClientFactory{queueBuilder: queueBuilder}
}

func (qcf *QueueClientFactory) Build(ctx context.Context) (IQueue, error) {
	return qcf.queueBuilder(ctx)
}
