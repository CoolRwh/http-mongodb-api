package db

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoDbClient *mongo.Client

func InitMongoDb() {
	url := viper.GetString("mongodb.url")
	minPoolSize := viper.GetUint64("mongodb.minPoolSize")
	maxPoolSize := viper.GetUint64("mongodb.maxPoolSize")
	// 设置MongoDB连接URI
	clientOptions := options.Client().SetMinPoolSize(minPoolSize).SetMaxPoolSize(maxPoolSize).ApplyURI(url)
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
	log.Println("init mongodb success!")
}

func CloseMongoDbClient() {
	err := MongoDbClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("MongoDB connection closed")
}
