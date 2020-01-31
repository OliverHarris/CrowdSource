package product

import (
	"context"
	"log"
	"reflect"
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Product information
type Product struct {
	ProductName pName
	Ingredients pIng
	Nutrition   pNutrition
	ID          string `bson:"_id"`
	Error       string
	Version     int64
}

type pName struct {
	Name    string
	Votes   PerVote
	Users   []UserVote
	Changes []pName
	Stamp   int64
	Vote    bool
}
type pIng struct {
	Ingredients []string
	Votes       PerVote
	Changes     []pIng
	Users       []UserVote
	Stamp       int64
	Vote        bool
}

type pNutrition struct {
	Nutrition   map[string][2]float32
	Weight      string
	Recommended string
	Votes       PerVote
	Changes     []pNutrition
	Users       []UserVote
	Stamp       int64
	Vote        bool
}

type PerVote struct {
	UpHigh   int
	UpLow    int
	DownHigh int
	DownLow  int
	Users    []UserVote
}

type UserVote struct {
	User string
	Up   bool
}

func GetProductInfo(barcode string, username string, conn *mongo.Database) Product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	var finalProduct Product
	err := doc.Decode(&finalProduct)
	if err != nil {
		log.Printf("error %s", err)
		finalProduct.Error = "Product not found"
	} else {

		vc := VoteCheck{"INGREDIENTS", finalProduct.ID, finalProduct.Ingredients.Stamp, username, conn}
		finalProduct.Ingredients.Vote = canVote(vc)
		vc.part = "NAME"
		vc.version = finalProduct.ProductName.Stamp
		finalProduct.ProductName.Vote = canVote(vc)

		vc.part = "NUTRITION"
		vc.version = finalProduct.Nutrition.Stamp
		finalProduct.Nutrition.Vote = canVote(vc)

	}
	return finalProduct
}

func AlterProduct(p Product, username string, conn *mongo.Database) {
	//decide how many points they should get
	prod := GetProductInfo(p.ID, username, conn)
	sec := time.Now().Unix()
	if len(p.Ingredients.Ingredients) > 0 && !testEq(p.Ingredients.Ingredients, prod.Ingredients.Ingredients) {
		prod.Ingredients = pIng{Ingredients: p.Ingredients.Ingredients}
		prod.Ingredients.Stamp = sec
		point := user.Point{p.ID, sec, "INGREDIENTS", 1, false, sec}
		user.AddPoint(point, username, conn)
	}
	if len(p.Nutrition.Nutrition) > 0 && reflect.DeepEqual(p.Nutrition.Nutrition, prod.Nutrition.Nutrition) {
		prod.Nutrition = pNutrition{Nutrition: p.Nutrition.Nutrition}
		prod.Nutrition.Stamp = sec
		point := user.Point{p.ID, sec, "NUTRITION", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if p.ProductName.Name != "" && p.ProductName.Name != prod.ProductName.Name {
		prod.ProductName = pName{Name: p.ProductName.Name}
		prod.ProductName.Stamp = sec
		point := user.Point{p.ID, sec, "NAME", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	prod.Version = sec
	//Now insert it into the database
	collection := conn.Collection("products")
	filter := bson.M{"_id": p.ID}
	collection.FindOneAndReplace(context.TODO(), filter, p, options.FindOneAndReplace().SetUpsert(true))
	//collection.InsertOne(context.TODO(), p)
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
