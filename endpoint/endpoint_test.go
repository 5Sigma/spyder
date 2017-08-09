package endpoint

import (
	"github.com/Jeffail/gabs"
	"testing"
)

func buildConfig() *gabs.Container {
	json := gabs.New()
	json.Set("http://localhost/api/endpoint", "url")
	json.Set("GET", "method")
	return json
}

func TestLoad(t *testing.T) {
	json := gabs.New()
	json.Set("http://localhost/api/endpoint", "url")
	_, err := LoadBytes([]byte(json.Bytes()))
	if err == nil {
		t.Errorf("Should return error for missing required fields")
	}
	json.Set("GET", "method")
	if _, err := LoadBytes([]byte(json.Bytes())); err != nil {
		t.Errorf("Config parsing error: %s", err.Error())
	}
}

func TestRequestUrl(t *testing.T) {
	json := buildConfig()
	params, _ := json.Object("data")
	params.Set("3", "option1")
	params.Set("4", "option2")
	ep, _ := LoadBytes(json.Bytes())
	expectedUrl := "http://localhost/api/endpoint?option1=3&option2=4"
	if ep.RequestURL() != expectedUrl {
		t.Errorf("Request URL missmatch:\nExpecting: %s\nReceived: %s", expectedUrl,
			ep.RequestURL())
	}
	ep.Method = "POST"
	if ep.RequestURL() != ep.Url {
		t.Errorf("Request URL missmatch:\nExpecting: %s\nReceived: %s", ep.Url,
			ep.RequestURL())
	}

}

func TestHeaders(t *testing.T) {
	config := buildConfig()
	config.Object("headers")
	config.Set("application/json", "headers", "Content-Type")
	ep, err := LoadBytes(config.Bytes())
	if err != nil {
		t.Fatalf("Error reading config: %s", err.Error())
	}

	headerMap := ep.Headers()
	headerValues := headerMap["Content-Type"]
	if len(headerValues) == 0 {
		t.Fatalf("No headers returned")
	}

	if headerValues[0] != "application/json" {
		t.Errorf("Header not stored or retrieved correctly")
	}
}

func TestRequestMethod(t *testing.T) {
	config := buildConfig()
	config.Set("get", "method")
	ep, err := LoadBytes(config.Bytes())
	if err != nil {
		t.Fatalf("Error reading config: %s", err.Error())
	}
	if ep.RequestMethod() != "GET" {
		eFormat := "RequestMethod() returned %s instead of %s"
		t.Errorf(eFormat, ep.RequestMethod(), "GET")
	}
}
