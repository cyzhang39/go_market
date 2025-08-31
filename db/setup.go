package db

import (
	"context"
	"fmt"
	"log"
	"time"

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