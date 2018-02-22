package endpoint

import (
	"github.com/5sigma/spyder/config"
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
	_, err := LoadBytes("", []byte(json.Bytes()))
	if err == nil {
		t.Errorf("Should return error for missing required fields")
	}
	json.Set("GET", "method")
	if _, err := LoadBytes("", []byte(json.Bytes())); err != nil {
		t.Errorf("Config parsing error: %s", err.Error())
	}
}

func TestRequestUrl(t *testing.T) {

	//GET request with params
	json := buildConfig()
	params, _ := json.Object("data")
	params.Set("3", "option1")
	params.Set("4", "option2")
	ep, _ := LoadBytes("", json.Bytes())
	expectedUrl := "http://localhost/api/endpoint?option1=3&option2=4"
	if ep.RequestURL() != expectedUrl {
		t.Errorf("Request URL missmatch:\nExpecting: %s\nReceived: %s", expectedUrl,
			ep.RequestURL())
	}

	// POST request
	json.Set("post", "method")
	ep, _ = LoadBytes("", json.Bytes())
	if ep.RequestURL() != ep.Url {
		t.Errorf("Request URL missmatch:\nExpecting: %s\nReceived: %s", ep.Url,
			ep.RequestURL())
	}

	// GET request with variable expansion
	config.LocalConfig.SetVariable("var", "value1")
	config.LocalConfig.SetVariable("host", "127.0.0.1")
	json = buildConfig()
	params, _ = json.Object("data")
	params.Set("$var", "option2")
	params.Set("3", "option1")
	json.Set("http://$host/api/endpoint", "url")
	ep, _ = LoadBytes("", json.Bytes())
	expectedUrl = "http://127.0.0.1/api/endpoint?option1=3&option2=value1"
	if ep.RequestURL() != expectedUrl {
		t.Errorf("Request URL missmatch:\nExpecting: %s\nReceived: %s", expectedUrl,
			ep.RequestURL())
	}

}

func TestHeaders(t *testing.T) {
	epConfig := buildConfig()
	epConfig.Object("headers")
	epConfig.Set("application/json", "headers", "Content-Type")
	epConfig.Set("$var", "headers", "x-custom")
	config.LocalConfig.SetVariable("var", "value1")
	ep, err := LoadBytes("", epConfig.Bytes())
	if err != nil {
		t.Fatalf("Error reading config: %s", err.Error())
	}

	headerMap := ep.Headers
	contentTypeValues := headerMap["Content-Type"]
	if contentTypeValues[0] != "application/json" {
		t.Errorf("Header not stored or retrieved correctly")
	}

	customValues := headerMap["x-custom"]
	if customValues[0] != "value1" {
		t.Errorf("Header not stored or retrieved correctly")
	}

}

func TestRequestMethod(t *testing.T) {
	epConfig := buildConfig()
	epConfig.Set("get", "method")
	ep, err := LoadBytes("", epConfig.Bytes())
	if err != nil {
		t.Fatalf("Error reading config: %s", err.Error())
	}
	if ep.RequestMethod() != "GET" {
		eFormat := "RequestMethod() returned %s instead of %s"
		t.Errorf(eFormat, ep.RequestMethod(), "GET")
	}
	config.LocalConfig.SetVariable("method", "POST")
	epConfig.Set("$method", "method")
	ep, err = LoadBytes("", epConfig.Bytes())
	if err != nil {
		t.Fatalf("Error reading config: %s", err.Error())
	}
	if ep.RequestMethod() != "POST" {
		t.Fatalf("String expansion for request method failed\n%s\n",
			ep.RequestMethod())
	}
}

func TestRequestData(t *testing.T) {
	config.LocalConfig.SetVariable("var", "123")
	configStr := `
		{
			"method": "POST",
			"url": "http://localhost:8081/api/cool",
			"data": {
				"dynamic": "$var"
			}
		}
	`
	ep, err := LoadBytes("", []byte(configStr))
	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}
	parsedData, _ := gabs.ParseJSON(ep.RequestData())
	v := parsedData.Path("dynamic").Data().(string)
	if v != "123" {
		t.Errorf("String expansion failed\n%s\n%s", v, "123")
	}
}
