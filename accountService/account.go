package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Username      string `bson:"_id"`
	Points        int
	PointsHistory []histpoints
}
type histpoints struct {
	Item      string
	Version   int
	Type      string
	Points    int
	Confirmed bool
	Timestamp int64
}

func getAccountDB(username string) (user, error) {
	var account user
	collection := conn.Collection("user")
	filter := bson.M{"_id": username}
	doc := collection.FindOne(context.TODO(), filter)
	err := doc.Decode(&account)
	return account, err
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	account, err := getAccountDB(getUsername(r))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	output, _ := json.Marshal(account)
	w.Write(output)

}

func getBoard(w http.ResponseWriter, r *http.Request) {
	scores := getTopScores(5)
	output, _ := json.Marshal(scores)
	w.Write(output)
}

type score struct {
	Username string `bson:"_id"`
	Points   int
}

func getTopScores(numberToGet int) []score {
	collection := conn.Collection("user")
	filter := bson.M{} //Could find with "public" tag
	findOptions := options.Find()
	findOptions.SetLimit(int64(numberToGet))
	findOptions.SetSort(bson.D{{"points", -1}})
	docs, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var scores []score
	docs.All(context.TODO(), &scores)
	return scores
}
