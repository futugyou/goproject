package domain

// IAggregate represents the basic interface for aggregates.
type IAggregate interface {
	AggregateName() string
	AggregateId() string
}
