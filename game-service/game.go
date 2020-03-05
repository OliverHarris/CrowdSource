package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type game struct {
	Session   string `bson:"_id"`
	Points    int
	Active    bool
	Questions []string
	Stamp     int64
}

type play struct {
	Barcode string
	Session string
}

type playResult struct {
	Correct bool
	Found   bool
	Error   string
}

func tokenGenerator() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func generateSession(w http.ResponseWriter, r *http.Request) {
	//Check to see if there is a running session
	collection := conn.Collection("game")
	filter := bson.M{"user": getUsername(r), "active": true}
	doc := collection.FindOne(context.TODO(), filter)
	var runningGame game
	err := doc.Decode(&runningGame)
	if err != nil {
		log.Printf("error %s", err)
		//Create a new token
		runningGame.Session = tokenGenerator()
		runningGame.Active = true
		runningGame.Stamp = time.Now().Unix()
		collection.InsertOne(context.TODO(), runningGame)
	}
	output, _ := json.Marshal(runningGame)
	w.Write(output)
}

func end(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p play
	err := decoder.Decode(&p)
	failOnError(err, "Failed to decode play")

	delete(sessions, p.Session)

	filter := bson.D{{"_id", p.Session}}

	update := bson.M{"active": false}
	collection := conn.Collection("game")
	collection.UpdateOne(context.TODO(), filter, update)

}

func games(w http.ResponseWriter, r *http.Request) {
	collection := conn.Collection("game")
	filter := bson.M{"user": getUsername(r)}
	cur, err := collection.Find(context.TODO(), filter)
	failOnError(err, "Failed to collect many")
	var results []*game
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem game
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	output, _ := json.Marshal(results)
	w.Write(output)
}

var sessions map[string]string

func playOne(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var p play
	err := decoder.Decode(&p)
	failOnError(err, "Failed to decode play")

	var re playResult
	if _, ok := sessions[p.Session]; !ok {
		//doesn't exist
		re.Error = "No game session"
		output, _ := json.Marshal(re)
		w.Write(output)
		return
	}

	prod := GetProductInfo(p.Barcode)

	if prod.Error == "Product not found" {
		re.Correct = false
		re.Found = false
	} else {
		//Our only question is if the product name is correct
		re.Found = true
		correct := sessions[p.Session]
		correct = correct[len(correct)-1:]
		if strings.ToLower(prod.ProductName.Name) == strings.ToLower(correct) {
			re.Correct = true
		} else {
			re.Correct = false
		}
	}

	delete(sessions, p.Session)
	incPoint(p.Session)
	output, _ := json.Marshal(re)
	w.Write(output)
}

func incPoint(session string) {
	filter := bson.D{{"_id", session}}

	update := bson.D{
		{"$inc", bson.D{
			{"points", 1},
		}},
	}
	collection := conn.Collection("game")
	collection.UpdateOne(context.TODO(), filter, update)

}
