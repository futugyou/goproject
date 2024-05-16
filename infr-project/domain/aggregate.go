package domain

// IAggregate represents the basic interface for aggregates.
type IAggregate interface {
	AggregateName() string
	AggregateId() string
}

type Aggregate struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (a *Aggregate) AggregateId() string {
	return a.Id
}
