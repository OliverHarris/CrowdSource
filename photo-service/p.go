package main

import (
	"bytes"
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

func clearString(st string) []string {
	st = strings.ToLower(st)
	st = strings.ReplaceAll(st, "(", ",")
	minAscii := 97
	maxAscii := 122
	var final bytes.Buffer
	for _, value := range st {
		ascii := int(value)
		if ascii >= minAscii && ascii <= maxAscii || ascii == 44 || ascii == 32 {
			final.Write([]byte(string(ascii)))

		}
	}

	newst := strings.Split(final.String(), ",")
	var finalSet []string
	index := make(map[string]string)
	for _, val := range newst {
		val = strings.TrimSpace(val)
		if val != "" {
			if _, ok := index[val]; !ok {
				index[val] = ""
				finalSet = append(finalSet, val)
			}
		}
	}

	return finalSet
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
	text2 := clearString(text)
	log.Printf("New text %s", text2)
	//split the string; should do this in a function to remove repeated etc and %
	currentProduct.Ingredients = text2

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
