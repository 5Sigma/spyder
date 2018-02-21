package endpoint

import (
	"errors"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/faker"
	"github.com/5sigma/spyder/sectionfile"
	"github.com/Jeffail/gabs"
	"net/url"
	"strings"
)

type (
	EndpointConfig struct {
		Name            string
		Filename        string
		json            *gabs.Container
		Method          string
		Url             string
		OnComplete      []string
		Transform       []string
		Headers         map[string][]string
		Prompts         []*Prompt
		ResponsePattern string
		RequestPattern  string
	}
	Prompt struct {
		Name         string
		DefaultValue string
	}
)

func New() *EndpointConfig {
	return &EndpointConfig{
		json:       &gabs.Container{},
		OnComplete: []string{},
		Transform:  []string{},
	}
}

// Load - Loads a confugruation from a file on the disk.
func Load(filename string) (*EndpointConfig, error) {
	sf, err := sectionfile.Load(filename)
	if err != nil {
		return nil, err
	}
	epConfig, err := LoadBytes(filename, []byte(sf.Contents("endpoint")))
	if err != nil {
		return nil, err
	}
	epConfig.ResponsePattern = sf.Contents("response")
	epConfig.RequestPattern = sf.Contents("request")
	return epConfig, nil
}

// LoadBytes - Loads a configuration from a byte array
func LoadBytes(filename string, fileBytes []byte) (*EndpointConfig, error) {
	var (
		jsonObject *gabs.Container
		err        error
		epConfig   *EndpointConfig
	)

	jsonObject, err = gabs.ParseJSON(fileBytes)

	if err != nil {
		return nil, err
	}
	if !validate(jsonObject) {
		return nil, errors.New("Invalid endpoint configuration")
	}
	method, _ := jsonObject.Path("method").Data().(string)
	url, _ := jsonObject.Path("url").Data().(string)

	headerMap := make(map[string][]string)
	children, _ := jsonObject.S("headers").ChildrenMap()
	for key, child := range children {
		headerMap[key] = []string{config.ExpandString(child.Data().(string))}
	}

	promptsJson, _ := jsonObject.S("prompts").ChildrenMap()
	prompts := []*Prompt{}
	for key, child := range promptsJson {
		var defaultValue = ""
		if child.Exists("defaultValue") {
			defaultValue = child.Path("defaultValue").Data().(string)
		}
		prompt := &Prompt{Name: key, DefaultValue: defaultValue}
		prompts = append(prompts, prompt)
	}
	name, _ := jsonObject.Path("name").Data().(string)

	epConfig = &EndpointConfig{
		Name:       name,
		json:       jsonObject,
		Method:     method,
		Url:        url,
		OnComplete: []string{},
		Transform:  []string{},
		Headers:    headerMap,
		Prompts:    prompts,
		Filename:   filename,
	}

	transformNodes, _ := jsonObject.S("transform").Children()
	for _, node := range transformNodes {
		epConfig.Transform = append(epConfig.Transform, node.Data().(string))
	}

	onCompleteNodes, _ := jsonObject.S("onComplete").Children()
	for _, node := range onCompleteNodes {
		epConfig.OnComplete = append(epConfig.OnComplete, node.Data().(string))
	}

	return epConfig, nil
}

// GetString - returns a string from an arbitrary path in the configuration.
func (ep *EndpointConfig) GetString(path string) string {
	str, _ := ep.json.Path(path).Data().(string)
	return config.ExpandString(str)
}

// GetJSONString - returns the inner JSON at the path as a string.
func (ep *EndpointConfig) GetJSONString(path string) string {
	if ep.json.Exists("data") {
		return ep.json.Path("data").String()
	}
	return ""
}

// GetJSONBytes - returns the inner JSON at the path as a byte array.
func (ep *EndpointConfig) GetJSONBytes(path string) []byte {
	return ep.json.Path("data").Bytes()
}

// RequestMethod - returns the request method.
func (ep *EndpointConfig) RequestMethod() string {
	method := strings.ToUpper(ep.GetString("method"))
	return method
}

// RequestURL - returns the full url for the request. If this is a GET request
// and has request parameters they are included in the URL.
func (ep *EndpointConfig) RequestURL() string {
	if ep.RequestMethod() == "GET" {
		baseURL, _ := url.Parse(faker.ExpandFakes(config.ExpandString(ep.Url)))
		params := url.Values{}
		for k, v := range ep.GetRequestParams() {
			params.Add(k, v)
		}
		baseURL.RawQuery = params.Encode()
		return config.ExpandURL(baseURL.String())
	} else {
		return config.ExpandURL(ep.GetString("url"))
	}
}

// GetRequestParams - Returns a string map of any request params for the
// request. This only applies to GET requests.
func (ep *EndpointConfig) GetRequestParams() map[string]string {
	if ep.RequestMethod() != "GET" {
		return make(map[string]string)
	}
	paramsMap := make(map[string]string)
	children, err := ep.json.S("data").ChildrenMap()
	if err != nil {
		return paramsMap
	}
	for key, child := range children {
		childData, ok := child.Data().(string)
		if ok {
			paramsMap[key] = faker.ExpandFakes(config.ExpandString(childData))
		}

	}
	return paramsMap
}

// RequestData - returns the data attribute from the config. This contains the
// payload, for a POST request, that will be sent to the server.
func (ep *EndpointConfig) RequestData() []byte {
	dataJSON := ep.GetJSONString("data")
	dataJSON = config.ExpandString(dataJSON)
	dataJSON = faker.ExpandFakes(dataJSON)
	return []byte(dataJSON)
}

func (ep *EndpointConfig) SetRequestData(data map[string]interface{}) error {
	ep.json.SetP(data, "data")
	return nil
}

// validate - Validates that the configuration is valid and has the required
// parameters.
func validate(json *gabs.Container) bool {
	if !json.ExistsP("method") {
		return false
	}
	if !json.ExistsP("url") {
		return false
	}
	return true
}

func (ep *EndpointConfig) String() string {
	return ep.json.String()
}

func (ep *EndpointConfig) ResponseDefinition() *gabs.Container {
	node := ep.json.S("definition", "response")
	return node
}

func (ep *EndpointConfig) RequestDefinition() *gabs.Container {
	node := ep.json.S("definition", "request")
	return node
}
