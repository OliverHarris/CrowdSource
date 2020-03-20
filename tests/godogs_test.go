package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

const base = "http://project:8900"

type apiFeature struct {
	response *resty.Response
	request  *resty.Request
	token    string
}

func (a *apiFeature) resetResponse(interface{}) {
	a.response = &resty.Response{}
	a.request = resty.New().R().SetAuthToken(a.token)
}

func getToken() string {

	url := "https://dev-4x8d9r5y.eu.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"61LADdan2MAmjGI0FG8JXpWLsghX3UsC\",\"client_secret\":\"urdA8PNFNeevVrvJypQFQrCvgQl72KMC2SjMyH4dvEaKOIsMuH1jsT8KesL7oczT\",\"audience\":\"https://project.gateway\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var t map[string]interface{}
	err := json.Unmarshal(body, &t)
	if err != nil {
		log.Printf("Error.. %s", err.Error())
	}
	token := fmt.Sprintf("%s", t["access_token"])
	log.Printf("Token %s", t["access_token"])
	return token
}

func (a *apiFeature) iPutThisInformation(arg1 *gherkin.DocString) (err error) {
	a.request.SetBody(arg1.Content)
	return
}

func (a *apiFeature) theResponseShouldBe(arg1 int) error {
	if arg1 != a.response.StatusCode() {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", arg1, a.response.StatusCode())
	}
	return nil

}

func (a *apiFeature) sendRequestTo(method, endpoint string) (err error) {
	end := fmt.Sprintf("%s%s", base, endpoint)
	switch method {
	case "POST":
		log.Printf("i'm sending %v", a.request.Body)

		a.response, err = a.request.Post(end)
		break
	case "GET":

		a.response, err = a.request.Get(end)
		break

	default:
		log.Printf("Help, i'm lost for the method")
	}
	return
}

func (a *apiFeature) theResponseShouldInclude(body *gherkin.DocString) (err error) {

	var expected, actual []byte
	var data interface{}
	// if err = json.Unmarshal([]byte(body.Content), &data); err != nil {
	// 	return
	// }
	if expected, err = json.Marshal(data); err != nil {
		return
	}
	actual = a.response.Body()
	if !bytes.Contains(actual, expected) {
		err = fmt.Errorf("The response %s doesn't contain %s", string(actual), string(body.Content))
	}
	return

}

func (a *apiFeature) theProductDoesExist(arg1 string) error {
	collection := conn.Collection("products")
	filter := bson.M{"_id": arg1}
	doc := collection.FindOne(context.TODO(), filter)
	if doc.Err() != nil {
		//Create the product
		//Be cheeky, use the rest api

		client := resty.New().R().SetAuthToken(a.token)
		client.SetBody("{\"productName\":{\"Test product\"},\"ID\":\"" + arg1 + "\", \"Ingredients\": {\"Ingredients\": [\"Wheat\",\"Egg\",\"Sugar\"]}")
		client.Post(fmt.Sprintf("%s%s%s", base, "/product/", arg1))

	}
	return nil
}

func (a *apiFeature) theProductDoesntExist(arg1 string) error {
	collection := conn.Collection("products")
	filter := bson.M{"_id": arg1}
	doc := collection.FindOne(context.TODO(), filter)
	if doc.Err() == nil {
		collection.DeleteOne(context.TODO(), filter)
	}
	return nil
}

var conn *mongo.Database

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}
	api.token = getToken()
	var err error
	//Create a database connection
	conn, err = configDB(context.Background())
	failOnError(err, "Connecting to database failed")

	s.BeforeScenario(api.resetResponse)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.sendRequestTo)
	s.Step(`^the response should be (\d+)$`, api.theResponseShouldBe)
	s.Step(`^the product "([^"]*)" does exist$`, api.theProductDoesExist)
	s.Step(`^the product "([^"]*)" doesn\'t exist$`, api.theProductDoesntExist)
	s.Step(`^I put this information$`, api.iPutThisInformation)
	s.Step(`^the response should include$`, api.theResponseShouldInclude)

}
