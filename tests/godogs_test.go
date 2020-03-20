package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

const base = "http://192.168.122.176:8900"

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

func (a *apiFeature) iPutThisInformation(arg1 *gherkin.DataTable) (err error) {
	m := make(map[string]interface{})

	for _, r := range arg1.Rows {
		first := ""
		for _, c := range r.Cells {
			if first == "" {
				first = c.Value
			} else {
				if strings.Contains(c.Value, "[") {
					var dat []string
					err := json.Unmarshal([]byte(c.Value), &dat)
					if err != nil {
						log.Printf("Error %v", err.Error())
					}
					m[first] = dat
				} else if strings.Contains(c.Value, "{") {
					var dat map[string]interface{}
					err := json.Unmarshal([]byte(c.Value), &dat)
					if err != nil {
						log.Printf("Error %v", err.Error())
					}
					m[first] = dat
				} else {
					m[first] = c.Value
				}

			}
		}
	}
	v, err := json.Marshal(m)
	a.request.SetBody(v)
	log.Printf("%v ", string(v))
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

func (a *apiFeature) theResponseShouldInclude(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func (a *apiFeature) theProductTestDoesExist(arg1 int) error {
	return godog.ErrPending
}

func (a *apiFeature) iRequestToCreateAProduct() error {
	return godog.ErrPending
}

func (a *apiFeature) theProductTestDoesntExist(arg1 int) error {
	return godog.ErrPending
}

func (a *apiFeature) theProductTestExists(arg1 int) error {
	return godog.ErrPending
}

func (a *apiFeature) iRequestTest(arg1 int) error {
	return godog.ErrPending
}

func (a *apiFeature) myResponseShouldInclude(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}
	api.token = getToken()

	s.BeforeScenario(api.resetResponse)
	s.Step(`^I put this information$`, api.iPutThisInformation)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.sendRequestTo)
	s.Step(`^the response should be (\d+)$`, api.theResponseShouldBe)
	s.Step(`^the response should include$`, api.theResponseShouldInclude)
	s.Step(`^the product test(\d+) does exist$`, api.theProductTestDoesExist)
	s.Step(`^I request to create a product$`, api.iRequestToCreateAProduct)
	s.Step(`^the product test(\d+) doesn\'t exist$`, api.theProductTestDoesntExist)
	s.Step(`^the product test(\d+) exists$`, api.theProductTestExists)
	s.Step(`^I request test(\d+)$`, api.iRequestTest)
	s.Step(`^my response should include$`, api.myResponseShouldInclude)
}
