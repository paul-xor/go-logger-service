package main

import (
	"context"
	"fmt"
	"log"
	"log-service/cmd/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "82"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Add your application logic here
	log.Println("Connected to MongoDB and application is running.")

	app := Config{
		Models: data.New(client),
	}
	// start webserver
	go func() {
		if err := app.serve(); err != nil {
			log.Panic(err)
		}
	}()
	select {}
}

func (app *Config) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting server on port %s\n", webPort)
	return srv.ListenAndServe()
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL).SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	// Ping the database to verify connection
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, err
	}

	return c, nil
}
