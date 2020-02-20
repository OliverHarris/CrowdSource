package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Settings is the users account settings
type Settings struct {
	Allergies            []string
	RecommendedNutrition map[string]float32
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	collection := conn.Collection("accountInfo")
	username := getUsername(r)
	filter := bson.M{"_id": username}
	doc := collection.FindOne(context.TODO(), filter)
	var account Settings
	err := doc.Decode(&account)
	if err != nil {
		log.Printf("error %s", err)
		//Account not found. This is fine, return an empty account
	}

	output, _ := json.Marshal(account)
	w.Write(output)
}

func setSettings(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var account Settings
	err := decoder.Decode(&account)
	failOnError(err, "Failed to decode product")
	collection := conn.Collection("accountInfo")
	username := getUsername(r)
	filter := bson.M{"_id": username}
	collection.FindOneAndReplace(context.TODO(), filter, account, options.FindOneAndReplace().SetUpsert(true))
}
