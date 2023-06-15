package db

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-talk/common/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func MongoInit() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongodbCfg.Address))
	if err != nil {
		logrus.Errorf("Connection MongoDB Error: %s", err.Error())
		return
	}
	Mongo = client.Database(config.MongodbCfg.DBName)
	logrus.Infof("Connected MongoDB success")
}
