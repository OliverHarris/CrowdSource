package product

import (
	"context"
	"reflect"
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Product information
type Product struct {
	ProductName pName
	Ingredients pIng
	Serving     pServing
	Nutrition   pNutrition
	ID          string `bson:"_id"`
	Error       string
	Version     int64
}

//New product from the client
type ProductInput struct {
	ID          string `bson:"_id"`
	ProductName string
	Ingredients []string
	Serving     string
	Nutrition   map[string]float32
}

type pName struct {
	Name    string
	Up      int
	Down    int
	Changes []pName
	Stamp   int64
}
type pIng struct {
	Ingredients []string
	Up          int
	Down        int
	Changes     []pIng
	Stamp       int64
}
type pServing struct {
	Serving string
	Up      int
	Down    int
	Changes []pServing
	Stamp   int64
}

type pNutrition struct {
	Nutrition map[string]float32
	Up        int
	Down      int
	Changes   []pNutrition
	Stamp     int64
}

func GetProductInfo(barcode string, conn *mongo.Database) Product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	var finalProduct Product
	err := doc.Decode(&finalProduct)
	if err != nil {
		finalProduct.Error = "Product not found"
	}
	return finalProduct
}

func AlterProduct(p ProductInput, username string, conn *mongo.Database) {
	//decide how many points they should get
	prod := GetProductInfo(p.ID, conn)
	sec := time.Now().Unix()
	if len(p.Ingredients) > 0 && !testEq(p.Ingredients, prod.Ingredients.Ingredients) {
		prod.Ingredients = pIng{Ingredients: p.Ingredients}
		prod.Ingredients.Stamp = sec
		point := user.Point{p.ID, sec, "INGREDIENTS", 1, false, sec}
		user.AddPoint(point, username, conn)
	}
	if len(p.Nutrition) > 0 && reflect.DeepEqual(p.Nutrition, prod.Nutrition.Nutrition) {
		prod.Nutrition = pNutrition{Nutrition: p.Nutrition}
		prod.Nutrition.Stamp = sec
		point := user.Point{p.ID, sec, "NUTRITION", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if p.ProductName != "" && p.ProductName != prod.ProductName.Name {
		prod.ProductName = pName{Name: p.ProductName}
		prod.ProductName.Stamp = sec
		point := user.Point{p.ID, sec, "NAME", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if p.Serving != "" && p.Serving != prod.Serving.Serving {
		prod.Serving = pServing{Serving: p.Serving}
		prod.Serving.Stamp = sec
		point := user.Point{p.ID, sec, "SERVING", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	prod.Version = sec
	//Now insert it into the database
	collection := conn.Collection("products")
	collection.InsertOne(context.TODO(), p)
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