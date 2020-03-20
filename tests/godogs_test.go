package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

const base = "http://192.168.122.176:8900"

type apiFeature struct {
	response *resty.Response
	request  *resty.Request
}

func (a *apiFeature) resetResponse(interface{}) {
	a.response = &resty.Response{}
	a.request = resty.New().R().SetAuthToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik5VWTVOek16TURSRlFURXhRa0k0TmtWRVFVSXlOa0UxTXpCRFF6WkNNVGM1TWtNMFF6Z3dOQSJ9.eyJpc3MiOiJodHRwczovL2Rldi00eDhkOXI1eS5ldS5hdXRoMC5jb20vIiwic3ViIjoiNjFMQURkYW4yTUFtakdJMEZHOEpYcFdMc2doWDNVc0NAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vcHJvamVjdC5nYXRld2F5IiwiaWF0IjoxNTg0MzgwNjA5LCJleHAiOjE1ODQ0NjcwMDksImF6cCI6IjYxTEFEZGFuMk1BbWpHSTBGRzhKWHBXTHNnaFgzVXNDIiwiZ3R5IjoiY2xpZW50LWNyZWRlbnRpYWxzIn0.fEPApsvBFUA_oKupv0gJBmFGfZ5yK6uwxZ_Zlh7KjbzveyQzbO_MH-6ykhF23D_vABeuaD3liUtveQWi6UVuZCHGdwtavFwgoQKhIcxr82FkbyciMRKYNbPKcPcSLos32HHTD7A9yAbbvj82WwJCZiRTHQeSSJQuJjfNvoBcR6CVDz2IDRGK8--yUKK32KyXlkUkThJVwkfZUAf9ML3ZnGJ7yrnwrRj_qzhF_B1n2yGglNj0e6H5Loy0cSwr19ghMMWxtg9UnCvHKIGsaBoXCnPzBNH_qHAA6MpTsPpld2f3ypN_ZysBqkb0SROLVzU1ntX2fDmDvi5tYWPpzXe8FA")
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

func (a *apiFeature) theUserHasSetSomeInformation() error {
	return godog.ErrPending
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

	s.BeforeScenario(api.resetResponse)
	s.Step(`^I put this information$`, api.iPutThisInformation)
	s.Step(`^send "([^"]*)" request to "([^"]*)"$`, api.sendRequestTo)
	s.Step(`^the response should be (\d+)$`, api.theResponseShouldBe)
	s.Step(`^The user has set some information$`, api.theUserHasSetSomeInformation)
	s.Step(`^the response should include$`, api.theResponseShouldInclude)
	s.Step(`^the product test(\d+) does exist$`, api.theProductTestDoesExist)
	s.Step(`^I request to create a product$`, api.iRequestToCreateAProduct)
	s.Step(`^the product test(\d+) doesn\'t exist$`, api.theProductTestDoesntExist)
	s.Step(`^the product test(\d+) exists$`, api.theProductTestExists)
	s.Step(`^I request test(\d+)$`, api.iRequestTest)
	s.Step(`^my response should include$`, api.myResponseShouldInclude)
}
