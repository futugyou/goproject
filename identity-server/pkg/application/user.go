package application

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/futugyou/domaincore/application"
	domaincore "github.com/futugyou/domaincore/domain"

	"github.com/futugyousuzu/identity-server/pkg/domain"
	"github.com/futugyousuzu/identity-server/pkg/viewmodel"
)

type UserService struct {
	innerService *application.AppService
	repository   domain.UserRepository
	emailService EmailService
}

func NewUserService(
	repository domain.UserRepository,
	unitOfWork domaincore.UnitOfWork,
	emailService EmailService) *UserService {
	return &UserService{
		repository:   repository,
		innerService: application.NewAppService(unitOfWork),
		emailService: emailService,
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

func (u *UserService) CreateUser(ctx context.Context, request viewmodel.CreateUserRequest) (*viewmodel.CreateUserResponse, error) {
	if len(request.Email) == 0 || len(request.Password) == 0 {
		return nil, fmt.Errorf("email or password is empty")
	}

	err := u.checkUserExist(ctx, request)
	if err != nil {
		return nil, err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	user := domain.NewUser(request.Name, request.Email, string(hashed))

	err = u.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		// TODO: send verify email
		return u.repository.Insert(ctx, *user)
	})
	if err != nil {
		return nil, err
	}

	return &viewmodel.CreateUserResponse{
		ID: user.ID,
	}, err
}

func (u *UserService) checkUserExist(ctx context.Context, request viewmodel.CreateUserRequest) error {
	var orConditions []domaincore.FilterExpr = []domaincore.FilterExpr{
		domaincore.Eq{
			Field: "email",
			Value: request.Email,
		},
	}

	if len(request.Name) > 0 {
		orConditions = append(orConditions, domaincore.Eq{
			Field: "name",
			Value: request.Name,
		})
	}

	var filter domaincore.FilterExpr
	if len(orConditions) > 1 {
		filter = domaincore.Or(orConditions)
	} else {
		filter = orConditions[0]
	}

	query := domaincore.NewQueryOptions(nil, nil, nil, filter)
	datas, err := u.repository.Find(ctx, query)

	if err != nil {
		return err
	}

	if len(datas) > 0 {
		return fmt.Errorf("user exist")
	}

	return nil
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
