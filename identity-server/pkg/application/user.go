package application

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/futugyou/domaincore/application"
	domaincore "github.com/futugyou/domaincore/domain"

	domain "github.com/futugyousuzu/identity-server/pkg/domain/user"
	"github.com/futugyousuzu/identity-server/pkg/dto"
	"github.com/futugyousuzu/identity-server/pkg/options"
	"github.com/futugyousuzu/identity-server/pkg/viewmodel"
)

type UserService struct {
	innerService  *application.AppService
	repository    domain.UserRepository
	logrepository domain.UserLoginRepository
	emailService  EmailService
	opts          *options.Options
}

func NewUserService(
	unitOfWork domaincore.UnitOfWork,
	repository domain.UserRepository,
	logrepository domain.UserLoginRepository,
	emailService EmailService,
	opts *options.Options,
) *UserService {
	return &UserService{
		repository:    repository,
		logrepository: logrepository,
		innerService:  application.NewAppService(unitOfWork),
		emailService:  emailService,
		opts:          opts,
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
	user := domain.NewUser(request.Name, request.Email, string(hashed), nil)

	err = u.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		// TODO: generate a self-verifiable JWT token
		// or other type of token, but it needs to be persisted.
		emailDto := u.createEmailDto(user.Email, "")
		u.emailService.SendEmail(ctx, emailDto)
		return u.repository.Insert(ctx, *user)
	})
	if err != nil {
		return nil, err
	}

	return &viewmodel.CreateUserResponse{
		ID: user.ID,
	}, err
}

func (s *UserService) Login(ctx context.Context, name, password string) (*viewmodel.LoginResponse, error) {
	filter := domaincore.Eq{
		Field: "name",
		Value: name,
	}

	query := domaincore.NewQueryOptions(nil, nil, nil, filter)
	datas, err := s.repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("user not exist")
	}

	user := datas[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("password not match")
	}

	userLogin := domain.NewUserLogin(user.ID)

	if err = s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.logrepository.Insert(ctx, *userLogin)
	}); err != nil {
		return nil, err
	}

	return &viewmodel.LoginResponse{
		ID:      userLogin.ID,
		UserID:  userLogin.UserID,
		LoginAt: time.Unix(userLogin.Timestamp, 0),
	}, err
}

func (s *UserService) createEmailDto(email string, token string) dto.EmailDTO {
	url := fmt.Sprintf("%s/verify/%s?token=%s", s.opts.EmailVerifyUrl, email, token)

	return dto.EmailDTO{
		From:    s.opts.EmailFromAddress,
		To:      email,
		Subject: "Activate your account",
		Text:    fmt.Sprintf("Hello,\n\nPlease activate your account by clicking the following link: %s\n\nThank you!", url),
		Html:    fmt.Sprintf("<p>Hello,</p><p>Please activate your account by clicking the following link: <a href='%s'>Activate Account</a></p><p>Thank you!</p>", url),
	}
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
