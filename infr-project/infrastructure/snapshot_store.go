package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error)
	SaveSnapshot(ctx context.Context, aggregate EventSourcing) error
}

type ISnapshotStoreAsync[EventSourcing domain.IEventSourcing] interface {
	//  // LoadSnapshotAsync mothed usage
	//  func Stream(ctx context.Context, id string) error {
	//  	resultChan, errorChan := LoadSnapshotAsync(ctx, id)
	//  	select {
	//  	case datas := <-resultChan:
	//  	// handle data
	//  	case err := <-errorChan:
	//  	// handle error
	//  	case <-ctx.Done():
	//  	// handle timeout
	//  	}
	//  }
	//
	LoadSnapshotAsync(ctx context.Context, id string) (<-chan []EventSourcing, <-chan error)
	//  // SaveSnapshotAsync mothed usage
	//  func Stream(ctx context.Context, aggregate EventSourcing) error {
	//  	errorChan := SaveSnapshotAsync(ctx, aggregate)
	//  	select {
	//  	case err := <-errorChan:
	//  	// handle error
	//  	case <-ctx.Done():
	//  	// handle timeout
	//  	}
	//  }
	//
	SaveSnapshotAsync(ctx context.Context, aggregate EventSourcing) <-chan error
}
