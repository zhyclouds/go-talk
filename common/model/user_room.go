package model

import (
	"context"
	"go-talk/common/db"
	"go-talk/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// UserRoom 用户房间表, 存储在MongoDB中
type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"` // 房间 类型 【1-独聊房间 2-群聊房间】
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

// GetUserRoomByUserIdentityRoomIdentity 根据用户身份标识和房间身份标识查询用户房间
func GetUserRoomByUserIdentityRoomIdentity(userIdentity, roomIdentity string) (*UserRoom, error) {
	ur := new(UserRoom)
	err := db.Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"user_identity", userIdentity}, {"room_identity", roomIdentity}}).
		Decode(ur)
	return ur, err
}

// GetUserRoomByRoomIdentity 根据房间身份标识查询用户房间
func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	cursor, err := db.Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

// JudgeUserIsFriend 判断两个用户是否为好友
func JudgeUserIsFriend(userIdentity1, userIdentity2 string) bool {
	// 查询 userIdentity1 单聊房间列表
	cursor, err := db.Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{"user_identity", userIdentity1}, {"room_type", 1}})
	roomIdentities := make([]string, 0)
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return false
	}
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(ur)
		if err != nil {
			log.Logger.Error("Decode Error", zap.Error(err))
			return false
		}
		roomIdentities = append(roomIdentities, ur.RoomIdentity)
	}
	// 获取关联 userIdentity2 单聊房间个数
	cnt, err := db.Mongo.Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(), bson.M{"user_identity": userIdentity2, "room_type": 1, "room_identity": bson.M{"$in": roomIdentities}})
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return false
	}
	if cnt > 0 {
		return true
	}

	return false
}

// InsertOneUserRoom 插入一条用户房间记录, 建立好友关系
func InsertOneUserRoom(ur *UserRoom) error {
	_, err := db.Mongo.Collection(UserRoom{}.CollectionName()).InsertOne(context.Background(), ur)
	return err
}

// DeleteUserRoom 删除一条用户房间记录, 删除好友关系
func DeleteUserRoom(userIdentity, roomIdentity string) error {
	_, err := db.Mongo.Collection(UserRoom{}.CollectionName()).
		DeleteOne(context.Background(), bson.M{"user_identity": userIdentity, "room_identity": roomIdentity})
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return err
	}
	return nil
}

// GetUserRoomIdentity 获取两个用户的单聊房间身份标识
func GetUserRoomIdentity(userIdentity1, userIdentity2 string) string {
	// 查询 userIdentity1 单聊房间列表
	cursor, err := db.Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{"user_identity", userIdentity1}, {"room_type", 1}})
	roomIdentities := make([]string, 0)
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return ""
	}
	for cursor.Next(context.Background()) {
		ur := new(UserRoom)
		err := cursor.Decode(ur)
		if err != nil {
			log.Logger.Error("Decode Error", zap.Error(err))
			return ""
		}
		roomIdentities = append(roomIdentities, ur.RoomIdentity)
	}
	// 获取关联 userIdentity2 单聊房间个数
	ur := new(UserRoom)
	err = db.Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.M{"user_identity": userIdentity2, "room_type": 1, "room_identity": bson.M{"$in": roomIdentities}}).Decode(ur)
	if err != nil {
		log.Logger.Error("[DB ERROR]", zap.Error(err))
		return ""
	}

	return ur.RoomIdentity
}
