package application

import (
	"context"
	"fmt"

	coreapp "github.com/futugyou/domaincore/application"
	coredomain "github.com/futugyou/domaincore/domain"
	coreinfr "github.com/futugyou/domaincore/infrastructure"

	"github.com/futugyou/resourceservice/domain"
	models "github.com/futugyou/resourceservice/viewmodel"
)

type ResourceService struct {
	service    *coreapp.ApplicationService[domain.ResourceEvent, *domain.Resource]
	unitOfWork coredomain.UnitOfWork
}

func needStoreSnapshot(aggregate *domain.Resource) bool {
	return aggregate.AggregateVersion()%5 == 1
}

func NewResourceService(
	eventStore coreinfr.EventStore[domain.ResourceEvent],
	snapshotStore coreinfr.SnapshotStore[*domain.Resource],
	unitOfWork coredomain.UnitOfWork,
	eventPublisher coreinfr.EventDispatcher,
) *ResourceService {
	return &ResourceService{
		service:    coreapp.NewApplicationService(eventStore, snapshotStore, unitOfWork, domain.ResourceFactory, needStoreSnapshot, eventPublisher),
		unitOfWork: unitOfWork,
	}
}

func (s *ResourceService) CreateResource(ctx context.Context, aux models.CreateResourceRequest) (*models.CreateResourceResponse, error) {
	var res *domain.Resource
	resourceType := domain.GetResourceType(aux.Type)
	if err := s.service.WithUnitOfWork(ctx, func(ctx context.Context) error {
		res = domain.NewResource(aux.Name, resourceType, aux.Data, aux.ImageData, aux.Tags)
		return s.service.SaveSnapshotAndEvent(ctx, res)
	}); err != nil {
		return nil, err
	}

	return &models.CreateResourceResponse{ID: res.ID}, nil
}

func (s *ResourceService) UpdateResource(ctx context.Context, id string, aux models.UpdateResourceRequest) error {
	res, err := s.service.RetrieveLatestVersion(ctx, id)
	if err != nil {
		return err
	}

	source := *res
	oldVersion := source.Version
	aggregate, err := source.ChangeResource(aux.Name, source.Type, aux.Data, aux.ImageData, aux.Tags)
	if err != nil {
		return err
	}

	if aggregate == nil || oldVersion == aggregate.Version {
		return fmt.Errorf("the data in the resource has not changed")
	}

	return s.service.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}

// show all versions
func (s *ResourceService) AllVersionResource(ctx context.Context, id string) ([]models.ResourceView, error) {
	re, err := s.service.RetrieveAllVersions(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]models.ResourceView, len(re))
	for i := range re {
		result[i] = *convertResourceEntityToViewModel(re[i])
	}
	return result, nil
}

// show resource data for the specified version
func (s *ResourceService) ResourceWithVersion(ctx context.Context, id string, version int) (*models.ResourceView, error) {
	re, err := s.service.RetrieveSpecificVersion(ctx, id, version)
	if err != nil {
		return nil, err
	}

	return convertResourceEntityToViewModel(*re), nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id string) error {
	res, err := s.service.RetrieveLatestVersion(ctx, id)
	if err != nil {
		return err
	}

	if res == nil || (*res).ID == "" {
		return fmt.Errorf("resource: %s not found", id)
	}

	aggregate, err := (*res).DeleteResource()
	if err != nil {
		return err
	}

	return s.service.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}

func convertResourceEntityToViewModel(src *domain.Resource) *models.ResourceView {
	if src == nil {
		return nil
	}

	return &models.ResourceView{
		ID:        src.ID,
		Name:      src.Name,
		Type:      src.Type.String(),
		Data:      src.Data,
		ImageData: src.ImageData,
		Version:   src.Version,
		IsDelete:  src.IsDeleted,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		Tags:      src.Tags,
	}
}
