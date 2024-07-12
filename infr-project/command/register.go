package command

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/futugyou/infr-project/domain"
)

type Router struct {
	CommandBus *cqrs.CommandBus
}

func GetCommandHandler() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		BookRoomHandler{},
	}
}

func CreateCommandRouter() (*Router, error) {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	cqrsMarshaler := cqrs.JSONMarshaler{}
	commandBus, err := cqrs.NewCommandBusWithConfig(pubSub, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return params.CommandName, nil
		},
		Marshaler: cqrsMarshaler,
	})
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{}, nil)
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(middleware.Recoverer)
	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.CommandName, nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return pubSub, nil
			},
			Marshaler: cqrsMarshaler,
		},
	)
	if err != nil {
		return nil, err
	}

	handlers := GetCommandHandler()
	if err = commandProcessor.AddHandlers(handlers...); err != nil {
		return nil, err
	}

	go func() {
		router.Run(context.Background())
	}()

	<-router.Running()

	return &Router{
		CommandBus: commandBus,
	}, nil
}

type CommonHandler struct {
	unitOfWork domain.IUnitOfWork
}

func (s *CommonHandler) withUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := s.unitOfWork.Start(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.unitOfWork.Rollback(ctx)
		} else {
			commitErr := s.unitOfWork.Commit(ctx)
			if commitErr != nil {
				err = commitErr
			}
		}
		s.unitOfWork.End(ctx)
	}()

	err = fn(ctx)
	return err
}
