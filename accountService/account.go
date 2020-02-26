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

type user struct {
	Points Points
}

// Points is for the user (inside points.)
type Points struct {
	Scan    int
	Updates int `json:"-"`
	Deny    int `json:"-"`
	Trust   int
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

	//if the RDA is empty, we can set some defaults
	if len(account.RecommendedNutrition) == 0 {
		account.RecommendedNutrition = make(map[string]float32)
		account.RecommendedNutrition["Energy"] = 2000 // Energy is kcal
		account.RecommendedNutrition["Fat"] = 70
		account.RecommendedNutrition["Saturates"] = 20
		account.RecommendedNutrition["Carbohydrate"] = 260
		account.RecommendedNutrition["Sugar"] = 90
		account.RecommendedNutrition["Protein"] = 50
		account.RecommendedNutrition["Salt"] = 6
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

func getPoints(w http.ResponseWriter, r *http.Request) {
	collection := conn.Collection("user")
	username := getUsername(r)
	filter := bson.M{"_id": username}
	doc := collection.FindOne(context.TODO(), filter)
	var user user
	err := doc.Decode(&user)
	if err != nil {
		log.Printf("error %s", err)
		//No account
	}

	user.Points.Trust = getLevel(user.Points)

	output, _ := json.Marshal(user.Points)
	w.Write(output)
}

//copy from api/user
func getLevel(user Points) int {

	//Ignore the deny points for now
	if user.Updates > 5 && user.Scan > 5 {
		return 2
	} else if user.Scan > 5 {
		return 1
	} else {
		return 0
	}

}
