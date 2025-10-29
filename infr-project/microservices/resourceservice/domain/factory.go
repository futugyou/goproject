package domain

import (
	"reflect"

	"github.com/futugyou/domaincore/domain"
)

func ResourceFactory() *Resource {
	return &Resource{}
}

type ResourceEventProcessor interface {
	EventApply(r *Resource, event domain.DomainEvent)
}

var eventNameMappingEventTypes = map[string]reflect.Type{
	"ResourceCreated": reflect.TypeOf(&ResourceCreatedEvent{}),
	"ResourceUpdated": reflect.TypeOf(&ResourceUpdatedEvent{}),
	"ResourceDeleted": reflect.TypeOf(&ResourceDeletedEvent{}),
}

var ResourceEventProcessors = map[reflect.Type]ResourceEventProcessor{
	reflect.TypeOf(&ResourceCreatedEvent{}): &ResourceCreatedProcessor{},
	reflect.TypeOf(&ResourceUpdatedEvent{}): &ResourceUpdatedProcessor{},
	reflect.TypeOf(&ResourceDeletedEvent{}): &ResourceDeletedProcessor{},
}

type ResourceCreatedProcessor struct{}
type ResourceUpdatedProcessor struct{}
type ResourceDeletedProcessor struct{}

func (p *ResourceCreatedProcessor) EventApply(r *Resource, event domain.DomainEvent) {
	e := event.(*ResourceCreatedEvent)
	r.ID = e.ID
	r.Version = e.Version()
	r.CreatedAt = e.CreatedAt
	r.Name = e.Name
	r.Type = GetResourceType(e.Type)
	r.Data = e.Data
	r.ImageData = e.ImageData
	r.Tags = e.Tags
}

func (p *ResourceUpdatedProcessor) EventApply(r *Resource, event domain.DomainEvent) {
	e := event.(*ResourceUpdatedEvent)
	r.ID = e.ID
	r.Name = e.Name
	r.Type = GetResourceType(e.Type)
	r.Version = e.Version()
	r.Data = e.Data
	r.Tags = e.Tags
	r.ImageData = e.ImageData
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceDeletedProcessor) EventApply(r *Resource, event domain.DomainEvent) {
	e := event.(*ResourceDeletedEvent)
	r.ID = e.ID
	r.Version = e.Version()
	r.IsDeleted = true
	r.UpdatedAt = e.CreatedAt
}
