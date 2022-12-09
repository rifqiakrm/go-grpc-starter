package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

const (
	exampleURL        = "http://localhost:8081"
	userLoginEndpoint = exampleURL + "/v1/auth/login"
)

var (
	ctx    = context.Background()
	client = http.DefaultClient

	httpStatus int
	httpBody   []byte
)

func TestMain(_ *testing.M) {
	status := godog.TestSuite{
		Name:                "user v1",
		ScenarioInitializer: InitializeScenario,
	}.Run()

	os.Exit(status)
}

func restoreDefaultState(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(restoreDefaultState)

	ctx.Step(`^I login with username "([^"]*)" and password "([^"]*)"$`, iLoginWithUsernameAndPassword)
	ctx.Step(`^response status code must be (\d+)$`, responseStatusCodeMustBe)
	ctx.Step(`^response must match json:$`, responseMustMatchJSON)
}

func iLoginWithUsernameAndPassword(username, password string) error {
	body := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
	return callEndpoint(http.MethodPost, userLoginEndpoint, body)
}

func responseStatusCodeMustBe(code int) error {
	if httpStatus != code {
		return fmt.Errorf("expected HTTP status code %d, but got %d", code, httpStatus)
	}
	return nil
}

func responseMustMatchJSON(want *godog.DocString) error {
	return deepCompareJSON([]byte(want.Content), httpBody)
}

func deepCompareJSON(want, have []byte) error {
	var expected interface{}
	var actual interface{}

	err := json.Unmarshal(want, &expected)
	if err != nil {
		return err
	}
	err = json.Unmarshal(have, &actual)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func callEndpoint(method, url string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatus = resp.StatusCode
	httpBody, err = ioutil.ReadAll(resp.Body)
	return err
}

func updateResource(url string, requests *godog.Table) error {
	return mutateResource(http.MethodPut, url, requests)
}

func createResource(url string, requests *godog.Table) error {
	return mutateResource(http.MethodPost, url, requests)
}

func mutateResource(method, url string, requests *godog.Table) error {
	for _, row := range requests.Rows {
		body := strings.NewReader(row.Cells[0].Value)
		if err := callEndpoint(method, url, body); err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
