package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var conn *mongo.Database

func main() {
	//load the env file
	err := godotenv.Load()
	failOnError(err, "Error getting env vars")

	//Create a database connection
	conn, err = configDB(context.Background())
	failOnError(err, "Connecting to database failed")

	router := mux.NewRouter()

	router.HandleFunc("/product/{barcode}", getProduct).Methods("GET")
	router.HandleFunc("/product/{barcode}", updateProduct).Methods("POST")
	router.HandleFunc("/vote/{barcode}", voteOnProduct).Methods("POST")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func getUsername(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	log.Printf("Length %v", len(authHeaderParts))
	if len(authHeaderParts) < 2 {
		log.Printf("Token not found. Giving test username")
		return "test"
	}
	tokenString := authHeaderParts[1]
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)
	if err != nil {
		log.Printf("Token not found. Giving test username (error) %e", err)
		return "test"
	}
	claims, _ := token.Claims.(*CustomClaims)
	fmt.Printf("(user) %s ", claims.Subject)
	return claims.Subject

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
