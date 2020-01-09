package product

import (
	"context"
	"time"

	"../user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Product information
type Product struct {
	ProductName string
	Ingredients []string
	Serving     string
	Nutrition   map[string]float32
	Version     int
	ID          string `bson:"_id"`
	Trust       map[string]points
	Changed     string
	Error       string
	Changes     []Product
}
type points struct {
	User    string
	Confirm int
	Deny    int
	Version int
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

func CreateProduct(p Product, username string, conn *mongo.Database) {
	//decide how many points they should get
	if len(p.Ingredients) > 0 {
		point := user.Point{p.ID, p.Version, "INGREDIENTS", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if len(p.Nutrition) > 0 {
		point := user.Point{p.ID, p.Version, "NUTRITION", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if p.ProductName != "" {
		point := user.Point{p.ID, p.Version, "NAME", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	if p.Serving != "" {
		point := user.Point{p.ID, p.Version, "SERVING", 1, false, time.Now().Unix()}
		user.AddPoint(point, username, conn)
	}
	//Add the trust system
	po := points{username, 0, 0, 0}
	p.Trust["ProductName"] = po
	p.Trust["Ingredients"] = po
	p.Trust["Serving"] = po
	p.Trust["Nutrition"] = po

	//Now insert it into the database
	collection := conn.Collection("products")
	collection.InsertOne(context.TODO(), p)
}
