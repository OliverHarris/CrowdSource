package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/otiai10/gosseract"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type product struct {
	ProductName string   `bson:"product_name"`
	Ingredients []string `bson:"ingredients"`
	Serving     string
	Nutrition   map[string]float32
	Version     int
	ID          string `bson:"_id"`
	Trust       map[string]points
	Changed     string
	Error       string
	Changes     []product
}

type productImg struct {
	ProductName string `bson:"product_name"`
	Ingredients string `bson:"ingredients"`
	Serving     string
	Nutrition   string
	Version     int
	ID          string `bson:"_id"`
	Trust       map[string]points
	Changed     string
	Error       string
	Changes     []product
}
type points struct {
	User    string
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
func updateTheProduct(p productImg, barcode string) {
	currentProduct := getFromCrowd(barcode)

	client := gosseract.NewClient()
	defer client.Close()
	sDec, err := b64.StdEncoding.DecodeString(p.Ingredients)
	if err != nil {
		log.Fatal(err)
	}
	client.SetImageFromBytes(sDec)

	text, _ := client.Text()
	log.Printf("New text %s", text)
	//split the string; should do this in a function to remove repeated etc and %
	currentProduct.Ingredients = strings.Split(text, " ")

	//update the database
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	collection.FindOneAndReplace(context.TODO(), filter, currentProduct, options.FindOneAndReplace().SetUpsert(true))

}

func updateProduct(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := mux.Vars(r)
	barcode := params["barcode"]
	var product productImg
	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	updateTheProduct(product, barcode)
}
