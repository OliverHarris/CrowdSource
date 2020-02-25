package product

import (
	"context"
	"log"
	"reflect"
	"strconv"
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
	Scans       []string `json:"-"`
}

type pName struct {
	Name    string
	Votes   PerVote
	Users   []UserVote `json:"-"`
	Changes []pName
	Stamp   int64
	Vote    bool
}
type pIng struct {
	Ingredients []string
	Votes       PerVote
	Changes     []pIng
	Users       []UserVote `json:"-"`
	Stamp       int64
	Vote        bool
}

type pNutrition struct {
	Nutrition   map[string][2]float32
	Weight      string
	Recommended string
	Votes       PerVote
	Changes     []pNutrition
	Users       []UserVote `json:"-"`
	Stamp       int64
	Vote        bool
}

type PerVote struct {
	UpHigh    int `json:"-"`
	UpLow     int `json:"-"`
	DownHigh  int `json:"-"`
	DownLow   int `json:"-"`
	TrustUp   int
	TrustDown int
}

type UserVote struct {
	User string
	Up   bool
}

func GetProductInfo(barcode string, username string, conn *mongo.Database) *Product {
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

	finalProduct.Ingredients.Votes.TrustUp, finalProduct.Ingredients.Votes.TrustDown = Trust(finalProduct.Ingredients.Votes)
	finalProduct.ProductName.Votes.TrustUp, finalProduct.ProductName.Votes.TrustDown = Trust(finalProduct.ProductName.Votes)
	finalProduct.Nutrition.Votes.TrustUp, finalProduct.Nutrition.Votes.TrustDown = Trust(finalProduct.Nutrition.Votes)

	return &finalProduct
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
		p.Ingredients.Stamp = sec
		prod.Ingredients.Changes = nil
		p.Ingredients.Changes = append(prod.Ingredients.Changes, prod.Ingredients)
		prod.Ingredients.Users = nil
		p.Ingredients.Users = append(prod.Ingredients.Users, UserVote{username, true})
		if len(p.Ingredients.Changes) > 3 {
			p.Ingredients.Changes = p.Ingredients.Changes[1:]
		}
		p.Ingredients.Votes.DownHigh = 0
		p.Ingredients.Votes.DownLow = 0
		switch level {
		case 0:
			p.Ingredients.Votes.UpHigh = 0
			p.Ingredients.Votes.UpLow = 1
		default:
			p.Ingredients.Votes.UpHigh = 1
			p.Ingredients.Votes.UpLow = 0
		}
		c := pIng{p.Ingredients.Ingredients, p.Ingredients.Votes,
			p.Ingredients.Changes, p.Ingredients.Users, p.Ingredients.Stamp, false}

		prod.Ingredients = c
	}
	if reflect.DeepEqual(p.Nutrition.Nutrition, prod.Nutrition.Nutrition) ||
		p.Nutrition.Recommended != prod.Nutrition.Recommended ||
		p.Nutrition.Weight != prod.Nutrition.Weight {
		p.Nutrition.Stamp = sec
		prod.Nutrition.Changes = nil
		p.Nutrition.Changes = append(prod.Nutrition.Changes, prod.Nutrition)
		prod.Nutrition.Users = nil
		p.Nutrition.Users = append(prod.Nutrition.Users, UserVote{username, true})
		if len(p.Nutrition.Changes) > 3 {
			p.Nutrition.Changes = p.Nutrition.Changes[1:]
		}
		p.Nutrition.Votes.DownHigh = 0
		p.Nutrition.Votes.DownLow = 0
		switch level {
		case 0:
			p.Nutrition.Votes.UpHigh = 0
			p.Nutrition.Votes.UpLow = 1
		default:
			p.Nutrition.Votes.UpHigh = 1
			p.Nutrition.Votes.UpLow = 0
		}
		c := pNutrition{p.Nutrition.Nutrition, p.Nutrition.Weight,
			p.Nutrition.Recommended, p.Nutrition.Votes,
			p.Nutrition.Changes, p.Nutrition.Users, p.Nutrition.Stamp, false}

		prod.Nutrition = calcRecommended(c)
	}
	if p.ProductName.Name != "" && p.ProductName.Name != prod.ProductName.Name {

		p.ProductName.Stamp = sec
		prod.ProductName.Changes = nil
		p.ProductName.Changes = append(prod.ProductName.Changes, prod.ProductName)
		prod.ProductName.Users = nil
		p.ProductName.Users = append(prod.ProductName.Users, UserVote{username, true})
		if len(p.ProductName.Changes) > 3 {
			p.ProductName.Changes = p.ProductName.Changes[1:]
		}
		p.ProductName.Votes.DownHigh = 0
		p.ProductName.Votes.DownLow = 0
		switch level {
		case 0:
			p.ProductName.Votes.UpHigh = 0
			p.ProductName.Votes.UpLow = 1
		default:
			p.ProductName.Votes.UpHigh = 1
			p.ProductName.Votes.UpLow = 0
		}
		c := pName{p.ProductName.Name, p.Ingredients.Votes, p.ProductName.Users,
			p.ProductName.Changes, p.ProductName.Stamp, false}

		prod.ProductName = c

	}
	log.Printf("Users2.. %v", prod.ProductName.Users)
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

func calcRecommended(info pNutrition) pNutrition {
	weight, _ := strconv.ParseFloat(info.Weight, 32)
	recommended, _ := strconv.ParseFloat(info.Recommended, 32)
	if weight < 1 && recommended < 1 {
		return info
	}

	for k, v := range info.Nutrition {
		//v[0] is the value for the total
		oneGram := v[0] / float32(weight)
		recommendedGram := oneGram * float32(recommended)
		//Set the recommended value
		temp := info.Nutrition[k]
		temp[1] = recommendedGram
		info.Nutrition[k] = temp
	}

	return info
}
