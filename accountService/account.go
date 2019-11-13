package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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
