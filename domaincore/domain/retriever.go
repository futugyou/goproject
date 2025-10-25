package domain

type AggregateRetriever[Aggregate AggregateRoot] interface {
	RetrieveAllVersions(id string) ([]Aggregate, error)
	RetrieveSpecificVersion(id string, version int) (*Aggregate, error)
	RetrieveLatestVersion(id string) (*Aggregate, error)
}
