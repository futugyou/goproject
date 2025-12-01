package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	domaincore "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyousuzu/identity-server/pkg/domain/user"
)

type UserRepository struct {
	mongoimpl.BaseCRUD[user.User]
}

func NewUserRepository(client *mongo.Client, config mongoimpl.DBConfig) *UserRepository {
	if config.CollectionName == "" {
		config.CollectionName = "users"
	}

	getID := func(a user.User) string { return a.AggregateId() }

	return &UserRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return u.checkNameOrEmail(ctx, email, "email")
}

func (u *UserRepository) FindByName(ctx context.Context, name string) (*user.User, error) {
	return u.checkNameOrEmail(ctx, name, "name")
}

func (s *UserRepository) checkNameOrEmail(ctx context.Context, str, field string) (*user.User, error) {
	filter := domaincore.Eq{
		Field: field,
		Value: str,
	}
	condition := domaincore.NewQueryOptions(nil, nil, nil, filter)
	users, err := s.Find(ctx, condition)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

var _ user.UserRepository = new(UserRepository)
