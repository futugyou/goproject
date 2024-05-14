package domain

type IAggregateRetriever[Aggregate IAggregate] interface {
	RetrieveAllVersions(id string) ([]Aggregate, error)
	RetrieveSpecificVersion(id string, version int) (*Aggregate, error)
	RetrieveLatestVersion(id string) (*Aggregate, error)
}
