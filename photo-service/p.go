package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/otiai10/gosseract"
)

type productImg struct {
	Ingredients string `bson:"ingredients"`
}

func clearString(st string) []string {
	st = strings.ToLower(st)
	st = strings.ReplaceAll(st, "(", ",")
	minAscii := 97
	maxAscii := 122
	var final bytes.Buffer
	for _, value := range st {
		ascii := int(value)
		if ascii >= minAscii && ascii <= maxAscii || ascii == 44 || ascii == 32 {
			final.Write([]byte(string(ascii)))

		}
	}

	newst := strings.Split(final.String(), ",")
	var finalSet []string
	index := make(map[string]string)
	for _, val := range newst {
		val = strings.TrimSpace(val)
		if val != "" {
			if _, ok := index[val]; !ok {
				index[val] = ""
				finalSet = append(finalSet, val)
			}
		}
	}

	return finalSet
}

func readIngText(p productImg) productImg {
	client := gosseract.NewClient()
	defer client.Close()
	sDec, err := b64.StdEncoding.DecodeString(p.Ingredients)
	if err != nil {
		log.Fatal(err)
	}
	client.SetImageFromBytes(sDec)

	text, _ := client.Text()
	log.Printf("New text %s", text)
	text2 := clearString(text)
	log.Printf("New text %s", text2)
	p.Ingredients = text
	return p

}

func readIng(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var product productImg
	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	text := readIngText(product)
	output, _ := json.Marshal(text)
	w.Write(output)

}
