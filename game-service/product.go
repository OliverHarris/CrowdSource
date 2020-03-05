package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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
	Nutrition   map[string][2]float64
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

func GetProductInfo(barcode string) *Product {
	collection := conn.Collection("products")
	filter := bson.M{"_id": barcode}
	doc := collection.FindOne(context.TODO(), filter)
	var finalProduct Product
	err := doc.Decode(&finalProduct)
	if err != nil {
		log.Printf("error %s", err)
		finalProduct.Error = "Product not found"
	}

	return &finalProduct
}
