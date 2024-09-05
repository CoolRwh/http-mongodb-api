package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoDbClient *mongo.Client

func InitMongoDb() {
	// 设置MongoDB连接URI
	clientOptions := options.Client().SetMinPoolSize(1).SetMaxPoolSize(10).ApplyURI("mongodb://root:jitao123@47.97.188.153:27018")
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	MongoDbClient = client
}
