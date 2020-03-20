package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func configDB(ctx context.Context) (*mongo.Database, error) {
	err := godotenv.Load()
	failOnError(err, "Error getting env vars")

	uri := fmt.Sprintf("mongodb://%s", os.Getenv("DBHOST"))
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	failOnError(err, "Couldn't connect to mongo")
	err = client.Connect(ctx)
	failOnError(err, "Mongo client couldn't connect with background context")
	todoDB := client.Database("pro")
	return todoDB, nil
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
