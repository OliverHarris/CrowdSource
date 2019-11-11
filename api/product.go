package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

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

func addProductData(barcode string, product product, user string) {
	currentProduct := getFromCrowd(barcode)
	product.Trust = make(map[string]points)
	product.ID = barcode
	changed := false
	p := points{user, 0, 0}
	if len(product.ProductName) != 0 {
		if currentProduct.ProductName == product.ProductName {
			p := points{currentProduct.Trust["ProductName"].User,
				currentProduct.Trust["ProductName"].Confirm + 1,
				currentProduct.Trust["ProductName"].Deny}
			product.Trust["ProductName"] = p
		} else {
			changed = true
			addPoints(1, false, user, barcode, "PRODUCTNAMEUPDATE")
			product.Trust["ProductName"] = p
		}
	}

	if len(product.Ingredients) != 0 {
		if testEq(product.Ingredients, currentProduct.Ingredients) {
			p := points{currentProduct.Trust["Ingredients"].User,
				currentProduct.Trust["Ingredients"].Confirm + 1, currentProduct.Trust["Ingredients"].Deny}
			product.Trust["Ingredients"] = p
		} else {
			changed = true
			addPoints(1, false, user, barcode, "INGREDIENTSUPDATE")
			product.Trust["Ingredients"] = p
		}
	}

	if len(product.Nutrition) != 0 {
		if reflect.DeepEqual(product.Nutrition, currentProduct.Nutrition) {
			p := points{currentProduct.Trust["Nutrition"].User,
				currentProduct.Trust["Nutrition"].Confirm + 1, currentProduct.Trust["Nutrition"].Deny}
			product.Trust["Nutrition"] = p
		} else {
			changed = true
			addPoints(1, false, user, barcode, "NUTRITIONUPDATE")
			product.Trust["Nutrition"] = p
		}
	}

	if len(product.Serving) != 0 {
		if product.Serving == currentProduct.Serving {
			p := points{currentProduct.Trust["Serving"].User,
				currentProduct.Trust["Serving"].Confirm + 1, currentProduct.Trust["Serving"].Deny}
			product.Trust["Serving"] = p
		} else {
			changed = true
			addPoints(1, false, user, barcode, "SERVINGUPDATE")
			product.Trust["Serving"] = p
		}
	}

	//if changed, keep old copy.
	if changed {
		product.Version = currentProduct.Version + 1
		product.Changes = append(currentProduct.Changes, currentProduct)
	} else {
		product.Changes = currentProduct.Changes
	}

	//update the database
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	collection.FindOneAndReplace(context.TODO(), filter, product, options.FindOneAndReplace().SetUpsert(true))

}

func getProductInfo(barcode string, user string) product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	var finalProduct product
	err := doc.Decode(&finalProduct)
	if err != nil {
		finalProduct.Error = "Product not found"
	}
	if !checkScanned(user, barcode) {
		addPoints(1, true, user, barcode, "SCAN")
	}
	return finalProduct
}

type user struct {
	user   string
	points int
}

type histpoints struct {
	Item      string
	Type      string
	Points    int
	Confirmed bool
	Timestamp int64
}

func checkScanned(username string, barcode string) bool {
	collection := conn.Collection("user")
	filter := bson.D{{"_id", username},
		{"pointsHistory.item", barcode}}
	doc, err := collection.Find(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	f := doc.Next(context.TODO())
	if !f {
		return false
	}
	return true

}

func addPoints(points int, confirmed bool, user string, barcode string, ptype string) {
	collection := conn.Collection("user")
	filter := bson.D{{"_id", user}}
	p := histpoints{barcode, ptype, points, confirmed, time.Now().Unix()}
	update := bson.D{
		{"$inc", bson.D{
			{"points", points},
		}},
		{"$set", bson.D{
			{"_id", user},
		}},
		{"$push", bson.D{
			{"pointsHistory", p},
		}},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updated a single document: ", updateResult.MatchedCount)

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
	product := getProductInfo(barcode, getUsername(r))
	if product.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
	product = getProductInfo(barcode, getUsername(r))
	if product.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	output, _ := json.Marshal(product)
	w.Write(output)

}

func testEq(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
