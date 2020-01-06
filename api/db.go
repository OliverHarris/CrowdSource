package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func configDB(ctx context.Context) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s", os.Getenv("DBHOST"))
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	failOnError(err, "Couldn't connect to mongo")
	err = client.Connect(ctx)
	failOnError(err, "Mongo client couldn't connect with background context")
	todoDB := client.Database("pro")
	return todoDB, nil
}
