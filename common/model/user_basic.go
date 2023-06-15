package model

import (
	"context"
	"go-talk/common/db"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Identity  string `bson:"identity"`
	Account   string `bson:"account"`
	Password  string `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int    `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

func (User) CollectionName() string {
	return "user_basic"
}

func GetUserBasicByAccountPassword(account, password string) (*User, error) {
	ub := new(User)
	err := db.Mongo.Collection(User{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"account", account}, {"password", password}}).
		Decode(ub)
	return ub, err
}

func GetUserBasicByIdentity(identity string) (*User, error) {
	ub := new(User)
	err := db.Mongo.Collection(User{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"identity", identity}}).
		Decode(ub)
	return ub, err
}

func GetUserBasicByAccount(account string) (*User, error) {
	ub := new(User)
	err := db.Mongo.Collection(User{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"account", account}}).
		Decode(ub)
	return ub, err
}

func GetUserBasicCountByEmail(email string) (int64, error) {
	return db.Mongo.Collection(User{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"email", email}})
}

func GetUserBasicCountByAccount(account string) (int64, error) {
	return db.Mongo.Collection(User{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"account", account}})
}

func InsertOneUserBasic(ub *User) error {
	_, err := db.Mongo.Collection(User{}.CollectionName()).InsertOne(context.Background(), ub)
	return err
}
