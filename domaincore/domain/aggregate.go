package domain

// IAggregate represents the basic interface for aggregates.
type AggregateRoot interface {
	AggregateId() string
}

type Aggregate struct {
	ID string `bson:"id" redis:"id" json:"id"`
}

func (a Aggregate) AggregateId() string {
	return a.ID
}
