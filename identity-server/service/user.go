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
}

func NewUserStore() *UserService {
	return &UserService{}
}

func (u *UserService) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
	entity, err := userRepo.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	entity.Password = ""
	return entity, nil
}

func (u *UserService) GetByName(ctx context.Context, name string) (*user.User, error) {
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
	entity, err := userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	entity.Password = ""
	return entity, nil
}

func (u *UserService) Login(ctx context.Context, name, password string) (*user.UserLogin, error) {
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
	userinfo, err := userRepo.FindByName(ctx, name)
	if err != nil {
		return nil, errors.New("user " + name + " can not find")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userinfo.Password), []byte(password))
	if err != nil {
		return nil, errors.New("user " + name + " password error")
	}

	userloginRepo := operator.UserLoginRepository
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
	operator := operate.DefaultOperator()
	userloginRepo := operator.UserLoginRepository
	return userloginRepo.Get(ctx, login_id)
}

func (u *UserService) CreateUser(ctx context.Context, userInfo user.User) error {
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
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
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
	entity, err := userRepo.FindByName(ctx, name)

	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	entity.Password = string(hashed)
	return userRepo.Update(ctx, entity, entity.ID)
}

func (u *UserService) ListUser(ctx context.Context) []*user.User {
	operator := operate.DefaultOperator()
	userRepo := operator.UserRepository
	result, _ := userRepo.GetAll(ctx)
	return result
}
