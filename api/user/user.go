package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"../product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type point struct {
	Item      string
	Version   int
	Type      string
	Points    int
	Confirmed bool
	Timestamp int64
}

// PointsForScan - Add point for scanning product if it hasn't been scanned
func PointsForScan(product product.Product, username string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username,
		"pointsHistory.item": product.ID}
	doc, _ := collection.Find(context.TODO(), filter, nil)
	if !doc.Next(context.TODO()) {
		//Add the points to the user
		p := point{product.ID, product.Version, "SCAN", 1, true, time.Now().Unix()}
		addPoint(p, username, conn)
	}
}

func addPoint(p point, user string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": user}
	update := bson.M{
		"$inc": bson.M{
			"points": p.Points,
		},
		"$set": bson.M{
			"_id": user,
		},
		"$push": bson.M{
			"pointsHistory": p,
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	failOnError(err, "Failed to add points")

	fmt.Println("updated a single document: ", updateResult.MatchedCount)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
