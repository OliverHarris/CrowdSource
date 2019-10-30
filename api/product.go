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
	ProductName string   `bson:"product_name"`
	Ingredients []string `bson:"ingredients"`
	Serving     string
	Nutrition   map[string]float32
	Version     int
	ID          string `bson:"_id"`
	Source      string
	Trust       map[string]points
	ChangeBy    changed
	Changed     string
	Error       string
}
type changed struct {
	Part string
	User string
}
type points struct {
	Confirm int
	Deny    int
}

func getFromCrowd(barcode string) product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	var finalProduct product
	err := doc.Decode(&finalProduct)
	if err != nil {
		finalProduct.Error = "Product not found"
	}
	return finalProduct
}

func addProductData(barcode string, product product, user string) {
	currentProduct := getFromCrowd(barcode)
	product.Trust = make(map[string]points)
	product.ID = barcode
	changed := false
	if len(product.ProductName) != 0 {
		if currentProduct.ProductName == product.ProductName {
			p := points{currentProduct.Trust["ProductName"].Confirm + 1, currentProduct.Trust["ProductName"].Deny}
			product.Trust["ProductName"] = p
		} else {
			changed = true
			p := points{0, 0}
			product.Trust["ProductName"] = p
		}
	}

	//if changed, keep old copy.
	if changed {

	}

	//update the database
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	collection.FindOneAndReplace(context.TODO(), filter, product, options.FindOneAndReplace().SetUpsert(true))

}

func getProductInfo(barcode string) product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	addPoint()
	var finalProduct product
	err := doc.Decode(&finalProduct)
	if err != nil {
		finalProduct.Error = "Product not found"
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

	product := getProductInfo(barcode)
	product.Source = "Open source database"
	output, _ := json.Marshal(product)
	w.Write(output)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := mux.Vars(r)
	barcode := params["barcode"]
	var product product
	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Decoded.. %v", product.ProductName)
	addProductData(barcode, product, "test")
}
