package model

import (
	"context"
	"go-talk/common/db"
	"go-talk/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type RoomBasic struct {
	Identity     string `bson:"identity"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}

func InsertOneRoomBasic(rb *RoomBasic) error {
	_, err := db.Mongo.Collection(RoomBasic{}.CollectionName()).InsertOne(context.Background(), rb)
	return err
}

func DeleteRoomBasic(roomIdentity string) error {
	_, err := db.Mongo.Collection(RoomBasic{}.CollectionName()).
		DeleteOne(context.Background(), bson.M{"identity": roomIdentity})
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return err
	}
	return nil
}
