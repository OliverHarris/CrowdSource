package main

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
)

func iPutThisInformation(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func theResponseShouldBe(arg1 int) error {
	return godog.ErrPending
}

func theUserHasSetSomeInformation() error {
	return godog.ErrPending
}

func theResponseShouldInclude(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func theProductTestDoesExist(arg1 int) error {
	return godog.ErrPending
}

func iRequestToCreateAProduct() error {
	return godog.ErrPending
}

func theProductTestDoesntExist(arg1 int) error {
	return godog.ErrPending
}

func theProductTestExists(arg1 int) error {
	return godog.ErrPending
}

func iRequestTest(arg1 int) error {
	return godog.ErrPending
}

func myResponseShouldInclude(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I put this information$`, iPutThisInformation)
	s.Step(`^the response should be (\d+)$`, theResponseShouldBe)
	s.Step(`^The user has set some information$`, theUserHasSetSomeInformation)
	s.Step(`^the response should include$`, theResponseShouldInclude)
	s.Step(`^the product test(\d+) does exist$`, theProductTestDoesExist)
	s.Step(`^I request to create a product$`, iRequestToCreateAProduct)
	s.Step(`^the product test(\d+) doesn\'t exist$`, theProductTestDoesntExist)
	s.Step(`^the product test(\d+) exists$`, theProductTestExists)
	s.Step(`^I request test(\d+)$`, iRequestTest)
	s.Step(`^my response should include$`, myResponseShouldInclude)
}
