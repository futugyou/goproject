package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/futugyousuzu/identity-server/operate"
	"github.com/futugyousuzu/identity-server/user"
)

type UserService struct {
	user.IUserRepository
	user.IUserLoginRepository
}

func NewUserService(operator *operate.Operator) *UserService {
	return &UserService{operator.UserRepository, operator.UserLoginRepository}
}

func (u *UserService) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	userRepo := u.IUserRepository
	entity, err := userRepo.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	entity.Password = ""
	return entity, nil
}

func (u *UserService) GetByName(ctx context.Context, name string) (*user.User, error) {
	userRepo := u.IUserRepository
	entity, err := userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	entity.Password = ""
	return entity, nil
}

func (u *UserService) Login(ctx context.Context, name, password string) (*user.UserLogin, error) {
	userRepo := u.IUserRepository
	userinfo, err := userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, errors.New("user " + name + " can not find")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userinfo.Password), []byte(password))
	if err != nil {
		return nil, errors.New("user " + name + " password error")
	}

	userloginRepo := u.IUserLoginRepository
	now := time.Now()
	hashed, _ := bcrypt.GenerateFromPassword([]byte(now.Format("20060102150405")+userinfo.ID), 14)
	userLogin := &user.UserLogin{
		ID:        string(hashed),
		UserID:    userinfo.ID,
		Timestamp: now.Unix(),
	}

	err = userloginRepo.Insert(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	return userLogin, nil
}

func (u *UserService) GetLoginInfo(ctx context.Context, login_id string) (*user.UserLogin, error) {
	userloginRepo := u.IUserLoginRepository
	return userloginRepo.Get(ctx, login_id)
}

func (u *UserService) CreateUser(ctx context.Context, userInfo user.User) error {
	userRepo := u.IUserRepository
	entity, _ := userRepo.FindByName(ctx, userInfo.Name)

	if len(entity.Name) != 0 {
		return errors.New("use exist")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 14)
	userInfo.Password = string(hashed)
	if len(userInfo.ID) == 0 {
		userInfo.ID = uuid.New().String()
	}

	return userRepo.Insert(ctx, &userInfo)
}

func (u *UserService) UpdatePassword(ctx context.Context, name, password string) error {
	userRepo := u.IUserRepository
	entity, err := userRepo.FindByName(ctx, name)

	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	entity.Password = string(hashed)
	return userRepo.Update(ctx, entity, entity.ID)
}

func (u *UserService) ListUser(ctx context.Context) []*user.User {
	userRepo := u.IUserRepository
	result, _ := userRepo.GetAll(ctx)
	return result
}
