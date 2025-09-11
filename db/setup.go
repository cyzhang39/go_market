package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDatabase() *mongo.Client {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:8080"))
	// client, err := mongo.Connect()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://testername:testerpass@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Failed to connect to mongodb")
		return nil
	}
	fmt.Println("Successfully connected to mongodb")
	return client
}

var Client *mongo.Client = MongoDatabase()

func CollectionDB(client *mongo.Client, name string) *mongo.Collection {
	return client.Database("goMarket").Collection(name)
}

var Chats *mongo.Collection
var Messages *mongo.Collection

func InitChats(client *mongo.Client, name string) error {
	Chats = client.Database(name).Collection("chats")
	Messages = client.Database(name).Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err := Chats.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "members", Value: 1}}, Options: options.Index().SetUnique(true)})
	if err != nil {
		log.Println("Chats create failed: ", err)
	}

	_, err = Messages.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "chatId", Value: 1}, {Key: "createdAt", Value: -1}}})
	if err != nil {
		log.Println("Messages create failed:", err)
	}

	return nil
}
