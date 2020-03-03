package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2/jwt"
)

var conn *mongo.Database
var debug bool

func main() {
	//load the env file
	err := godotenv.Load()
	failOnError(err, "Error getting env vars")

	debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}

	//Create a database connection
	conn, err = configDB(context.Background())
	failOnError(err, "Connecting to database failed")

	router := mux.NewRouter()

	router.HandleFunc("/product/{barcode}", getProduct).Methods("GET")
	router.HandleFunc("/product/{barcode}", changeProduct).Methods("POST")
	router.HandleFunc("/vote/{barcode}", productVote).Methods("POST")
	// router.HandleFunc("/vote/{barcode}", voteOnProduct).Methods("POST")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}

type CustomClaims struct {
	*jwt.Claims
	// additional claims apart from standard claims
	//We don't have any extra
	extra map[string]interface{}
}

func getUsername(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")

	if len(authHeaderParts) < 2 {
		log.Printf("Token not found. Giving test username")
		return "test"
	}
	tokenString := authHeaderParts[1]
	// token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)
	// if err != nil {
	// 	log.Printf("Token not found. Giving test username (error) %e", err)
	// 	return "test"
	// }
	// claims, _ := token.Claims.(*CustomClaims)#
	var claims CustomClaims
	// decode JWT token without verifying the signature
	token, _ := jwt.ParseSigned(tokenString)
	_ = token.UnsafeClaimsWithoutVerification(&claims)

	return claims.Subject

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
