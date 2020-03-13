package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/otiai10/gosseract"
)

type productImg struct {
	Ingredients string `bson:"ingredients"`
}

type productNutrition struct {
	Nutrition  string `bson:"ingredients"`
	Correction map[string]float64
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

func decodeNut(temptext string) map[string]float64 {
	st := strings.ToLower(temptext)
	total := make(map[string]float64)
	//break up the lines
	for _, line := range strings.Split(strings.TrimSuffix(st, "\n"), "\n") {
		log.Printf("Current line is %s", line)
		//if the line starts with "typical" ignore it
		if strings.Contains(line, "typical") {
			log.Printf("line started with typical")
			continue
		}

		//if the line starts with a number, ignore it
		remvspc := strings.Replace(line, " ", "", -1)
		onechar := []rune(remvspc)[0]

		if unicode.IsDigit(onechar) {
			log.Printf("This line starts with a digit %v", line)
			continue
		}

		p1 := ""
		p2 := ""
		if strings.Contains(line, "of which") {
			//if it contains a word, then a number then we can add it.
			cur := strings.Split(line, " ")
			p1 = fmt.Sprintf("%s %s %s", cur[0], cur[1], cur[2]) // Should be part 1
			t2 := cur[3]
			t2 = strings.Replace(t2, "g", "", -1)  //remove characters
			t2 = strings.Replace(t2, "kj", "", -1) //remove characters
			t2 = strings.Replace(t2, ",", ".", -1) //correct mistakes

			if _, err := strconv.ParseFloat(t2, 64); err == nil {
				fmt.Printf("looks like a number. %s\n", t2)
				p2 = t2
			} else {
				fmt.Printf("Not a number %s\n", t2)
				continue
			}

		} else {
			cur := strings.Split(line, " ")
			p1 = cur[0] // Should be part 1
			t2 := cur[1]
			t2 = strings.Replace(t2, "g", "", -1) //remove characters
			log.Printf("Remove g %s", t2)
			t2 = strings.Replace(t2, "kj", "", -1) //remove characters
			log.Printf("Remove kj %s", t2)
			t2 = strings.Replace(t2, ",", ".", -1) //correct mistakes

			if _, err := strconv.ParseFloat(t2, 64); err == nil {
				fmt.Printf("looks like a number. %s\n", t2)
				p2 = t2
			} else {
				fmt.Printf("Not a number %s\n", t2)
				continue
			}
		}
		total[p1], _ = strconv.ParseFloat(p2, 64)
		//fmt.Println(line)
	}

	return total
}

func readNutText(p productNutrition) productNutrition {
	client := gosseract.NewClient()
	defer client.Close()
	sDec, err := b64.StdEncoding.DecodeString(p.Nutrition)
	if err != nil {
		log.Fatal(err)
	}
	client.SetImageFromBytes(sDec)

	text, _ := client.Text()

	text2 := decodeNut(text)
	p.Nutrition = ""
	p.Correction = text2
	return p

}
func readNutrition(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var product productNutrition
	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	text := readNutText(product)
	output, _ := json.Marshal(text)
	w.Write(output)

}
