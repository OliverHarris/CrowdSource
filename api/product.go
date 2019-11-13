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
	Version int
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

func addProductData(barcode string, productIn product, user string) {
	currentProduct := getFromCrowd(barcode)
	productIn.Trust = make(map[string]points)
	productIn.ID = barcode
	changed := false
	p := points{user, 0, 0, 0}
	if len(productIn.ProductName) != 0 {
		if currentProduct.ProductName == productIn.ProductName {
			if canVote(barcode, user, "ProductNameVote+", currentProduct.Trust["ProductName"].Version) &&
				canVote(barcode, user, "ProductNameVote-", currentProduct.Trust["ProductName"].Version) {
				p := points{currentProduct.Trust["ProductName"].User,
					currentProduct.Trust["ProductName"].Confirm + 1,
					currentProduct.Trust["ProductName"].Deny,
					currentProduct.Trust["ProductName"].Version}
				productIn.Trust["ProductName"] = p
				addPoints(1, false, user, barcode, "ProductNameUpVote", currentProduct.Trust["ProductName"].Version)
			} else {
				productIn.Trust["ProductName"] = currentProduct.Trust["ProductName"]
			}

		} else {
			changed = true
			addPoints(1, false, user, barcode, "PRODUCTNAMEUPDATE", 0)
			p.Version = currentProduct.Trust["ProductName"].Version + 1
			productIn.Trust["ProductName"] = p
		}
	}

	if len(productIn.Ingredients) != 0 {
		if testEq(productIn.Ingredients, currentProduct.Ingredients) {
			if canVote(barcode, user, "IngredientsVote+", currentProduct.Trust["Ingredients"].Version) &&
				canVote(barcode, user, "IngredientsVote-", currentProduct.Trust["Ingredients"].Version) {
				p := points{currentProduct.Trust["Ingredients"].User,
					currentProduct.Trust["Ingredients"].Confirm + 1, currentProduct.Trust["Ingredients"].Deny,
					currentProduct.Trust["Ingredients"].Version}
				productIn.Trust["Ingredients"] = p
				addPoints(1, false, user, barcode, "ingredientsUpVote", currentProduct.Trust["Ingredients"].Version)
			} else {
				productIn.Trust["Ingredients"] = currentProduct.Trust["Ingredients"]
			}

		} else {
			changed = true
			addPoints(1, false, user, barcode, "INGREDIENTSUPDATE", 0)
			p.Version = currentProduct.Trust["Ingredients"].Version + 1
			productIn.Trust["Ingredients"] = p
		}
	}

	if len(productIn.Nutrition) != 0 {
		if reflect.DeepEqual(productIn.Nutrition, currentProduct.Nutrition) {

			if canVote(barcode, user, "NutritionVote+", currentProduct.Trust["Nutrition"].Version) &&
				canVote(barcode, user, "NutritionVote-", currentProduct.Trust["Nutrition"].Version) {
				p := points{currentProduct.Trust["Nutrition"].User,
					currentProduct.Trust["Nutrition"].Confirm + 1, currentProduct.Trust["Nutrition"].Deny,
					currentProduct.Trust["Nutrition"].Version}
				productIn.Trust["Nutrition"] = p
				addPoints(1, false, user, barcode, "NutritionUpVote", currentProduct.Trust["Nutrition"].Version)
			} else {
				productIn.Trust["Nutrition"] = currentProduct.Trust["Nutrition"]
			}

		} else {
			changed = true
			addPoints(1, false, user, barcode, "NUTRITIONUPDATE", 0)
			p.Version = currentProduct.Trust["Nutrition"].Version + 1
			productIn.Trust["Nutrition"] = p
		}
	}

	if len(productIn.Serving) != 0 {
		if productIn.Serving == currentProduct.Serving {
			if canVote(barcode, user, "ServingVote+", currentProduct.Trust["Serving"].Version) &&
				canVote(barcode, user, "ServingVote-", currentProduct.Trust["Serving"].Version) {
				p := points{currentProduct.Trust["Serving"].User,
					currentProduct.Trust["Serving"].Confirm + 1, currentProduct.Trust["Serving"].Deny,
					currentProduct.Trust["Serving"].Version}
				productIn.Trust["Serving"] = p
				addPoints(1, false, user, barcode, "ServingUpVote", currentProduct.Trust["Serving"].Version)
			} else {
				productIn.Trust["Serving"] = currentProduct.Trust["Serving"]
			}

		} else {
			changed = true
			addPoints(1, false, user, barcode, "SERVINGUPDATE", 0)
			p.Version = currentProduct.Trust["Serving"].Version + 1
			productIn.Trust["Serving"] = p
		}
	}

	//if changed, keep old copy.
	if changed {
		productIn.Version = currentProduct.Version + 1
		//Keep the last three changes
		if len(currentProduct.Changes) > 2 {
			productIn.Changes = append(currentProduct.Changes[1:], currentProduct)
		} else {
			productIn.Changes = append(currentProduct.Changes, currentProduct)
		}

	} else {
		productIn.Changes = currentProduct.Changes
	}

	//update the database
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	collection.FindOneAndReplace(context.TODO(), filter, productIn, options.FindOneAndReplace().SetUpsert(true))

}

func upVote(barcode string, username string, part string) {
	p := getProductInfo(barcode)

	if canVote(barcode, username, (part + "+"), p.Trust[part].Version) {
		//Upvote
		collection := conn.Collection("products")
		filter := bson.D{{"_id", barcode}}
		update := bson.D{
			{"$inc", bson.D{
				{"trust." + part + ".confirm", 1},
			}}}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("updated a single document: ", updateResult.MatchedCount)
	}
}

func downVote(barcode string, username string, part string) {
	p := getProductInfo(barcode)

	if canVote(barcode, username, (part + "+"), p.Trust[part].Version) {
		//Upvote
		collection := conn.Collection("products")
		filter := bson.D{{"_id", barcode}}
		update := bson.D{
			{"$inc", bson.D{
				{"trust." + part + ".deny", 1},
			}}}

		updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("updated a single document: ", updateResult.MatchedCount)
	}
}

type vote struct {
	part    string
	confirm bool
}

func voteOnProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]
	decoder := json.NewDecoder(r.Body)
	var voteJ vote
	err := decoder.Decode(&voteJ)
	if err != nil {
		log.Fatal(err)
	}
	if voteJ.confirm {
		upVote(barcode, getUsername(r), voteJ.part)
	} else {
		downVote(barcode, getUsername(r), voteJ.part)
	}
}

func canVote(barcode string, username string, part string, version int) bool {
	collection := conn.Collection("user")
	filter := bson.D{{"_id", username},
		{"pointsHistory.item", barcode},
		{"pointsHistory.type", part},
		{"pointsHistory.version", version}}
	doc, err := collection.Find(context.TODO(), filter, nil)
	if err != nil {
		log.Fatal(err)
	}
	f := doc.Next(context.TODO())
	if !f {
		log.Print("Can vote")
		return true
	}
	log.Print("Cant vote")
	return false
}

func getProductInfo(barcode string) product {
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

type user struct {
	user   string
	points int
}

type histpoints struct {
	Item      string
	Version   int
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

func addPoints(points int, confirmed bool, user string, barcode string, ptype string, version int) {
	collection := conn.Collection("user")
	filter := bson.D{{"_id", user}}
	p := histpoints{barcode, version, ptype, points, confirmed, time.Now().Unix()}
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

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]
	product := getProductInfo(barcode)
	if product.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !checkScanned(getUsername(r), barcode) {
		addPoints(1, true, getUsername(r), barcode, "SCAN", 0)
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
	product = getProductInfo(barcode)
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
