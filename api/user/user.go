package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Point is for the user
type Point struct {
	Item      string
	Version   int64
	Type      string
	Points    int
	Confirmed bool
	Timestamp int64
}

// PointsForScan - Add point for scanning product if it hasn't been scanned
func PointsForScan(barcode string, ver int64, username string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username,
		"pointsHistory.item": barcode}
	doc, _ := collection.Find(context.TODO(), filter, nil)
	if !doc.Next(context.TODO()) {
		//Add the points to the user
		p := Point{barcode, ver, "SCAN", 1, true, time.Now().Unix()}
		AddPoint(p, username, conn)
	}
}

//AddPoint
func AddPoint(p Point, user string, conn *mongo.Database) {
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
