package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error conneting:", err)
		return nil, err
	}

	return c, nil
}
