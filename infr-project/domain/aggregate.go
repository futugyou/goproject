package domain

// IAggregate represents the basic interface for aggregates.
type IAggregateRoot interface {
	AggregateName() string // this is use for storage table name, although ddd doesn’t need it
	AggregateId() string
}

type Aggregate struct {
	Id string `json:"id"`
}

func (a *Aggregate) AggregateId() string {
	return a.Id
}
