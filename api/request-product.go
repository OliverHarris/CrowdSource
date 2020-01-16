package main

import (
	"encoding/json"
	"net/http"

	"./product"
	"./user"
	"github.com/gorilla/mux"
)

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]
	product := product.GetProductInfo(barcode, conn)
	if product.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	username := getUsername(r)
	user.PointsForScan(product.ID, product.Version, username, conn)
	output, _ := json.Marshal(product)
	w.Write(output)

}

func changeProduct(w http.ResponseWriter, r *http.Request) {

}
