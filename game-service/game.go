package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type game struct {
	Session   string `bson:"_id"`
	Points    int
	Active    bool
	Questions []string
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
		collection.InsertOne(context.TODO(), runningGame)
	}
	output, _ := json.Marshal(runningGame)
	w.Write(output)
}
