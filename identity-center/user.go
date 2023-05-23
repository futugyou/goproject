package main

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	GetByName(ctx context.Context, name string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, name, password string) error
	ListUser(ctx context.Context) []User
}

type User struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type MongoUserSore struct {
	DBName         string
	CollectionName string
	client         *mongo.Client
}

func NewUserSore() *MongoUserSore {
	db := os.Getenv("db_name")
	c_name := "oauth2_users"
	url := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	return &MongoUserSore{
		DBName:         db,
		CollectionName: c_name,
		client:         client,
	}
}

func (u *MongoUserSore) GetByName(ctx context.Context, name string) (User, error) {
	c := u.client.Database(u.DBName).Collection(u.CollectionName)

	entity := new(User)
	filter := bson.D{{Key: "name", Value: name}}

	err := c.FindOne(ctx, filter).Decode(&entity)
	if err == nil {
		entity.Password = ""
	}
	return *entity, err
}

func (u *MongoUserSore) CreateUser(ctx context.Context, user User) error {
	c := u.client.Database(u.DBName).Collection(u.CollectionName)
	entity := new(User)
	filter := bson.D{{Key: "name", Value: user.Name}}
	c.FindOne(ctx, filter).Decode(&entity)

	if len(entity.Name) != 0 {
		return errors.New("use exist!")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashed)
	_, err := c.InsertOne(ctx, user)

	return err
}

func (u *MongoUserSore) UpdatePassword(ctx context.Context, name, password string) error {
	c := u.client.Database(u.DBName).Collection(u.CollectionName)
	entity := new(User)
	upsert := true
	option := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

	return c.FindOneAndUpdate(ctx, bson.D{{Key: "name", Value: name}}, bson.D{{Key: "password", Value: string(hashed)}}, &option).Decode(&entity)
}

func (u *MongoUserSore) ListUser(ctx context.Context) []User {
	result := make([]User, 0)
	coll := u.client.Database(u.DBName).Collection(u.CollectionName)
	filter := bson.D{}
	cursor, _ := coll.Find(ctx, filter)
	cursor.All(ctx, &result)
	for _, data := range result {
		cursor.Decode(&data)
		data.Password = ""
	}

	return result
}
