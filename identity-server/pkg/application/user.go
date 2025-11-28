package application

import (
	"context"
	"fmt"

	"github.com/futugyou/domaincore/application"
	domaincore "github.com/futugyou/domaincore/domain"

	"github.com/futugyousuzu/identity-server/pkg/domain"
	"github.com/futugyousuzu/identity-server/pkg/viewmodel"
)

type UserService struct {
	innerService *application.AppService
	repository   domain.UserRepository
}

func NewUserService(repository domain.UserRepository,
	unitOfWork domaincore.UnitOfWork) *UserService {
	return &UserService{
		repository:   repository,
		innerService: application.NewAppService(unitOfWork),
	}
}

func (s *UserService) SearchUser(ctx context.Context, request viewmodel.SearchUserRequest) ([]viewmodel.UserView, error) {
	var filter domaincore.FilterExpr = nil
	if request.Name != "" {
		filter = domaincore.Like{
			Field:           "name",
			Pattern:         request.Name,
			CaseInsensitive: true,
		}
	}

	query := domaincore.NewQueryOptions(nil, nil, nil, filter)
	datas, err := s.repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]viewmodel.UserView, len(datas))
	for i := range datas {
		result[i] = viewmodel.UserView{
			ID:            datas[i].ID,
			Name:          datas[i].Name,
			Email:         datas[i].Email,
			EmailVerified: datas[i].EmailVerified,
		}
	}

	return result, nil
}

func (s *UserService) CheckName(ctx context.Context, name string) error {
	return s.checkNameOrEmail(ctx, name, "name")
}

func (s *UserService) CheckEmail(ctx context.Context, email string) error {
	return s.checkNameOrEmail(ctx, email, "email")
}

func (s *UserService) checkNameOrEmail(ctx context.Context, str, field string) error {
	if len(str) == 0 {
		return fmt.Errorf("%s is empty", field)
	}

	filter := domaincore.Eq{
		Field: field,
		Value: str,
	}

	query := domaincore.NewQueryOptions(nil, nil, nil, filter)
	datas, err := s.repository.Find(ctx, query)
	if err != nil {
		return err
	}

	if len(datas) > 0 {
		return fmt.Errorf("%s exist", field)
	}

	return nil
}
