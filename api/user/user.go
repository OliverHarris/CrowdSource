package user

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Points Point
}

// Point is for the user (inside points.)
type Point struct {
	Scan    int
	Updates int
	Deny    int
}

func GetLevel(username string, conn *mongo.Database) int {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username}
	doc := collection.FindOne(context.TODO(), filter)
	var user user
	err := doc.Decode(&user)
	if err != nil {
		log.Println("User doesnt exist when getting level. Setting to 0")
		return 0
	}

	//Ignore the deny points for now
	if user.Points.Updates > 5 && user.Points.Scan > 5 {
		return 2
	} else if user.Points.Scan > 5 {
		return 1
	} else {
		return 0
	}

}

// PointsForScan - Add point for scanning product if it hasn't been scanned
func PointsForScan(username string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username}
	update := bson.M{
		"$inc": bson.M{
			"points.scan": 1,
		}, "$set": bson.M{
			"_id": username,
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	failOnError(err, "Failed to add points")

	fmt.Println("updated a single document: ", updateResult.MatchedCount)
}

func PointsForUpdate(username string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username}
	update := bson.M{
		"$inc": bson.M{
			"points.update": 1,
		}, "$set": bson.M{
			"_id": username,
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	failOnError(err, "Failed to add points")

	fmt.Println("updated a single document: ", updateResult.MatchedCount)
}

func PointsForDeny(username string, conn *mongo.Database) {
	collection := conn.Collection("user")
	filter := bson.M{"_id": username}
	update := bson.M{
		"$inc": bson.M{
			"points.deny": 1,
		}, "$set": bson.M{
			"_id": username,
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	failOnError(err, "Failed to add points")

	fmt.Println("updated a single document: ", updateResult.MatchedCount)
}

// //AddPoint
// func AddPoint(p Point, user string, conn *mongo.Database) {
// 	collection := conn.Collection("user")
// 	filter := bson.M{"_id": user}
// 	update := bson.M{
// 		"$inc": bson.M{
// 			"points": p.Points,
// 		},
// 		"$set": bson.M{
// 			"_id": user,
// 		},
// 		"$push": bson.M{
// 			"pointsHistory": p,
// 		},
// 	}
// 	updateResult, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
// 	failOnError(err, "Failed to add points")

// 	fmt.Println("updated a single document: ", updateResult.MatchedCount)

// }

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
