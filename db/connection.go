package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("URI_MONGO"))
	Client, _ = mongo.NewClient(clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err2 := Client.Connect(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}
	return Client
}

func GetDatabase() *mongo.Database {
	DB = Client.Database(os.Getenv("DATABASE"))
	return DB
}
