package resource

import (
	"reflect"

	domain "github.com/futugyou/infr-project/domain"
)

func ResourceFactory() *Resource {
	return &Resource{}
}

type ResourceEventProcessor interface {
	EventApply(r *Resource, event domain.IDomainEvent)
}

var eventNameMappingEventTypes = map[string]reflect.Type{
	"ResourceCreated":     reflect.TypeOf(&ResourceCreatedEvent{}),
	"ResourceUpdated":     reflect.TypeOf(&ResourceUpdatedEvent{}),
	"ResourceDeleted":     reflect.TypeOf(&ResourceDeletedEvent{}),
	"ResourceNameChanged": reflect.TypeOf(&ResourceNameChangedEvent{}),
	"ResourceDataChanged": reflect.TypeOf(&ResourceDataChangedEvent{}),
	"ResourceTagsChanged": reflect.TypeOf(&ResourceTagsChangedEvent{}),
	"ResourceTypeChanged": reflect.TypeOf(&ResourceTypeChangedEvent{}),
}

var ResourceEventProcessors = map[reflect.Type]ResourceEventProcessor{
	reflect.TypeOf(&ResourceCreatedEvent{}):     &ResourceCreatedProcessor{},
	reflect.TypeOf(&ResourceUpdatedEvent{}):     &ResourceUpdatedProcessor{},
	reflect.TypeOf(&ResourceNameChangedEvent{}): &ResourceNameChangedProcessor{},
	reflect.TypeOf(&ResourceDataChangedEvent{}): &ResourceDataChangedProcessor{},
	reflect.TypeOf(&ResourceTagsChangedEvent{}): &ResourceTagsChangedProcessor{},
	reflect.TypeOf(&ResourceTypeChangedEvent{}): &ResourceTypeChangedProcessor{},
	reflect.TypeOf(&ResourceDeletedEvent{}):     &ResourceDeletedProcessor{},
}

type ResourceCreatedProcessor struct{}
type ResourceUpdatedProcessor struct{}
type ResourceNameChangedProcessor struct{}
type ResourceDataChangedProcessor struct{}
type ResourceTagsChangedProcessor struct{}
type ResourceTypeChangedProcessor struct{}
type ResourceDeletedProcessor struct{}

func (p *ResourceCreatedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceCreatedEvent)
	r.Id = e.Id
	r.Version = e.Version()
	r.CreatedAt = e.CreatedAt
	r.Name = e.Name
	r.Type = GetResourceType(e.Type)
	r.Data = e.Data
	r.Tags = e.Tags
}

func (p *ResourceUpdatedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceUpdatedEvent)
	r.Id = e.Id
	r.Name = e.Name
	r.Type = GetResourceType(e.Type)
	r.Version = e.Version()
	r.Data = e.Data
	r.Tags = e.Tags
	r.ImageData = e.ImageData
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceNameChangedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceNameChangedEvent)
	r.Id = e.Id
	r.Version = e.Version()
	r.Name = e.Name
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceDataChangedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceDataChangedEvent)
	r.Id = e.Id
	r.Version = e.Version()
	r.Data = e.Data
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceTagsChangedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceTagsChangedEvent)
	r.Id = e.Id
	r.Version = e.Version()
	r.Tags = e.Tags
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceTypeChangedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceTypeChangedEvent)
	r.Id = e.Id
	r.Type = GetResourceType(e.Type)
	r.Version = e.Version()
	r.UpdatedAt = e.CreatedAt
}

func (p *ResourceDeletedProcessor) EventApply(r *Resource, event domain.IDomainEvent) {
	e := event.(*ResourceDeletedEvent)
	r.Id = e.Id
	r.Version = e.Version()
	r.IsDeleted = true
	r.UpdatedAt = e.CreatedAt
}
