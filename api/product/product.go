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
	Scans       []string
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
		finalProduct.Ingredients.Vote = canVote(finalProduct.Ingredients.Users, username)
		finalProduct.ProductName.Vote = canVote(finalProduct.ProductName.Users, username)
		finalProduct.Nutrition.Vote = canVote(finalProduct.Nutrition.Users, username)

	}
	return finalProduct
}

func AddScanPoint(p Product, username string, conn *mongo.Database) {
	if !stringInSlice(username, p.Scans) {
		//They haven't scanned before. Add a point
		log.Println("Adding point for a scan")
		user.PointsForScan(username, conn)
		collection := conn.Collection("products")
		filter := bson.M{"_id": p.ID}
		change := bson.M{"$push": bson.M{"Scans": username}}
		collection.UpdateOne(context.TODO(), filter, change)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AlterProduct(p Product, username string, conn *mongo.Database) {
	//decide how many points they should get
	prod := GetProductInfo(p.ID, username, conn)
	sec := time.Now().Unix()
	level := user.GetLevel(username, conn)
	if len(p.Ingredients.Ingredients) > 0 && !testEq(p.Ingredients.Ingredients, prod.Ingredients.Ingredients) {
		prod.Ingredients = pIng{Ingredients: p.Ingredients.Ingredients}
		prod.Ingredients.Stamp = sec

		switch level {
		case 0:
			prod.Ingredients.Votes.UpLow++
		default:
			prod.Ingredients.Votes.UpHigh++
		}
		prod.Ingredients.Users = append(prod.Ingredients.Users, UserVote{username, true})
	}
	if len(p.Nutrition.Nutrition) > 0 && reflect.DeepEqual(p.Nutrition.Nutrition, prod.Nutrition.Nutrition) {
		prod.Nutrition = pNutrition{Nutrition: p.Nutrition.Nutrition}
		prod.Nutrition.Stamp = sec
		switch level {
		case 0:
			prod.Nutrition.Votes.UpLow++
		default:
			prod.Nutrition.Votes.UpHigh++
		}
		prod.Nutrition.Users = append(prod.Nutrition.Users, UserVote{username, true})
	}
	if p.ProductName.Name != "" && p.ProductName.Name != prod.ProductName.Name {
		prod.ProductName = pName{Name: p.ProductName.Name}
		prod.ProductName.Stamp = sec
		switch level {
		case 0:
			prod.ProductName.Votes.UpLow++
		default:
			prod.ProductName.Votes.UpHigh++
		}
		prod.ProductName.Users = append(prod.ProductName.Users, UserVote{username, true})
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
