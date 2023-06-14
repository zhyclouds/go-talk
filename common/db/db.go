package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var MySQL *gorm.DB
var Mongo *mongo.Database

func Init() {
	MySQLInit()
	MongoInit()
}
