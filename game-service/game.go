package main

import (
	"encoding/json"
	"net/http"
)

type question struct {
	Question string
}

func getQuestion(w http.ResponseWriter, r *http.Request) {
	q := question{"Find a product that starts with S"}
	output, _ := json.Marshal(q)
	w.Write(output)
}
