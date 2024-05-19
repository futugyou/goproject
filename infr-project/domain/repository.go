package domain

// Deprecated: IEventSourcingRepository is deprecated, use ISnapshotStore and IEventStore.
type IEventSourcingRepository[EventSourcing IEventSourcing] interface {
	Load(id string) (*EventSourcing, error)
	LoadAll(id string) ([]EventSourcing, error)
	Save(aggregate EventSourcing) error
}

type IRepository[Aggregate IAggregateRoot] interface {
	Get(id string) (*Aggregate, error)
	Delete(id string) error
	Update(aggregate Aggregate) error
	Insert(aggregate Aggregate) error
}
