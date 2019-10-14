package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type product struct {
	ProductName string `bson:"product_name"`
	ID          string
	Brands      string
	Source      string
}

func getProductInfo(barcode string) product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	addPoint()
	var finalProduct product
	err := doc.Decode(&finalProduct)
	if err != nil {
		log.Printf("Not found %s", err)
	}
	return finalProduct
}

func productFromGod(barcode string) product {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://project:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	addPoint()
	var finalProduct product
	err = doc.Decode(&finalProduct)
	if err != nil {
		log.Printf("Not found %s", err)
	}
	return finalProduct
}

type user struct {
	user   string
	points int
}

func addPoint() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://project:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("Olivers").Collection("Points")
	filter := bson.D{{"user", "test"}}
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
		{"$set", bson.D{
			{"user", "test"},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updated a single document: ", updateResult.MatchedCount)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]

	product := productFromGod(barcode)
	product.Source = "Open source database"
	output, _ := json.Marshal(product)
	w.Write(output)

}
