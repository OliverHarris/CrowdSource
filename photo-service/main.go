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
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Create a database connection
	conn, err = configDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	//To allow other sources, enable cors
	//router.Use(cors)

	router.HandleFunc("/product/{barcode}", updateProduct).Methods("POST")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}

//To open the API to other sources (Browser ui) this will allow CORS
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		})
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
