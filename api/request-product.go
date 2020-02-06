package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./product"
	"github.com/gorilla/mux"
)

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]
	p := product.GetProductInfo(barcode, getUsername(r), conn)
	log.Printf("Barcode %s with an error of? %s", barcode, p.Error)
	if p.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	username := getUsername(r)
	product.AddScanPoint(p, username, conn)
	output, _ := json.Marshal(p)
	w.Write(output)

}

func changeProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p product.Product
	err := decoder.Decode(&p)
	failOnError(err, "Failed to decode product")
	product.AlterProduct(p, getUsername(r), conn)
}

func productVote(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p product.Vote
	err := decoder.Decode(&p)
	failOnError(err, "Failed to decode product")
	product.VoteOnProduct(p, getUsername(r), conn)
}
