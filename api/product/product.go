package product

import (
	"context"

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

func CreateProduct(p Product, user string) {
	//decide how many points they should get
	var points int
	if len(p.Ingredients) > 0 {
		points++
	}
	if len(p.Nutrition) > 0 {
		points++
	}
	if p.ProductName != "" {
		points++
	}
	if p.Serving != "" {
		points++
	}
}
