package user_mongo

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/futugyousuzu/identity-server/user"
)

type MongoUserStore struct {
	DBName                  string
	UserCollectionName      string
	UserLoginCollectionName string
	client                  *mongo.Client
}

func NewUserStore() *MongoUserStore {
	db := os.Getenv("db_name")
	u_name := "oauth2_users"
	l_name := "oauth2_login"
	url := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	return &MongoUserStore{
		DBName:                  db,
		UserCollectionName:      u_name,
		UserLoginCollectionName: l_name,
		client:                  client,
	}
}

func (u *MongoUserStore) GetByUID(ctx context.Context, uid string) (user.User, error) {
	c := u.client.Database(u.DBName).Collection(u.UserCollectionName)

	entity := new(user.User)
	filter := bson.D{{Key: "_id", Value: uid}}

	err := c.FindOne(ctx, filter).Decode(&entity)
	if err == nil {
		entity.Password = ""
	}
	return *entity, err
}

func (u *MongoUserStore) GetByName(ctx context.Context, name string) (user.User, error) {
	c := u.client.Database(u.DBName).Collection(u.UserCollectionName)

	entity := new(user.User)
	filter := bson.D{{Key: "name", Value: name}}

	err := c.FindOne(ctx, filter).Decode(&entity)
	if err == nil {
		entity.Password = ""
	}
	return *entity, err
}

func (u *MongoUserStore) Login(ctx context.Context, name, password string) (user.UserLogin, error) {
	c := u.client.Database(u.DBName).Collection(u.UserCollectionName)
	entity := new(user.User)
	userLogin := new(user.UserLogin)
	filter := bson.D{{Key: "name", Value: name}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		return *userLogin, errors.New("user " + name + " can not find")
	}

	err = bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(password))
	if err != nil {
		return *userLogin, err
	}

	c = u.client.Database(u.DBName).Collection(u.UserLoginCollectionName)
	now := time.Now()
	hashed, _ := bcrypt.GenerateFromPassword([]byte(now.Format("20060102150405")+entity.ID), 14)
	userLogin = &user.UserLogin{
		ID:        string(hashed),
		UserID:    entity.ID,
		Timestamp: now.Unix(),
	}

	_, err = c.InsertOne(ctx, *userLogin)

	return *userLogin, err
}

func (u *MongoUserStore) GetLoginInfo(ctx context.Context, login_id string) (user.UserLogin, error) {
	c := u.client.Database(u.DBName).Collection(u.UserLoginCollectionName)
	var login user.UserLogin
	err := c.FindOne(ctx, bson.D{{Key: "_id", Value: login_id}}).Decode(&login)
	return login, err
}

func (u *MongoUserStore) CreateUser(ctx context.Context, userInfo user.User) error {
	c := u.client.Database(u.DBName).Collection(u.UserCollectionName)
	entity := new(user.User)
	filter := bson.D{{Key: "name", Value: userInfo.Name}}
	c.FindOne(ctx, filter).Decode(&entity)

	if len(entity.Name) != 0 {
		return errors.New("use exist")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 14)
	userInfo.Password = string(hashed)
	if len(userInfo.ID) == 0 {
		userInfo.ID = uuid.New().String()
	}
	_, err := c.InsertOne(ctx, userInfo)

	return err
}

func (u *MongoUserStore) UpdatePassword(ctx context.Context, name, password string) error {
	c := u.client.Database(u.DBName).Collection(u.UserCollectionName)
	entity := new(user.User)
	upsert := true
	option := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	return c.FindOneAndUpdate(ctx, bson.D{{Key: "name", Value: name}}, bson.D{{Key: "password", Value: string(hashed)}}, &option).Decode(&entity)
}

func (u *MongoUserStore) ListUser(ctx context.Context) []user.User {
	result := make([]user.User, 0)
	coll := u.client.Database(u.DBName).Collection(u.UserCollectionName)
	filter := bson.D{}
	cursor, _ := coll.Find(ctx, filter)
	cursor.All(ctx, &result)
	for _, data := range result {
		cursor.Decode(&data)
		data.Password = ""
	}

	return result
}
