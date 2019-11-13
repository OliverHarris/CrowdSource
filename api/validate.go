package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func validatePoints(confirm bool, barcode string, version int, part string) {
	collection := conn.Collection("user")
	action := "-"
	if confirm {
		action = "+"
	}
	log.Printf("Looking for %s with the thing as %s, version of %s", barcode, (part + "Vote" + action), version)
	// filter := bson.D{{"pointsHistory.$.item", barcode},
	// 	{"pointsHistory.$.type", (part + "Vote" + action)},
	// 	{"pointsHistory.$.version", version},
	// 	{"pointsHistory.$.confirmed", false}}
	filter := bson.D{{"pointsHistory", bson.M{
		"$elemMatch": bson.M{"item": barcode,
			"type":      (part + "Vote" + action),
			"version":   version,
			"confirmed": false}}}}
	update := bson.D{
		{"$set", bson.D{
			{"pointsHistory.$.confirmed", true},
		}}}

	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated %v", updateResult.MatchedCount)
}
