package endpoint

import (
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"strings"
)

type (
	EndpointConfig struct {
		Filename   string
		json       *gabs.Container
		Method     string
		Url        string
		OnComplete []string
		Transform  []string
	}
)

func Load(filename string) (*EndpointConfig, error) {
	var (
		fileBytes  []byte
		jsonObject *gabs.Container
		err        error
		epConfig   *EndpointConfig
	)

	fileBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jsonObject, err = gabs.ParseJSON(fileBytes)
	if err != nil {
		return nil, err
	}

	method, _ := jsonObject.Path("method").Data().(string)
	url, _ := jsonObject.Path("url").Data().(string)

	epConfig = &EndpointConfig{
		Filename:   filename,
		json:       jsonObject,
		Method:     method,
		Url:        url,
		OnComplete: []string{},
		Transform:  []string{},
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

func (ep *EndpointConfig) GetString(path string) string {
	str, _ := ep.json.Path(path).Data().(string)
	return str
}

func (ep *EndpointConfig) GetJSONString(path string) string {
	return ep.json.Path("data").String()
}

func (ep *EndpointConfig) GetJSONBytes(path string) []byte {
	return ep.json.Path("data").Bytes()
}

func (ep *EndpointConfig) Headers() map[string][]string {
	headerMap := make(map[string][]string)
	children, _ := ep.json.S("headers").ChildrenMap()
	for key, child := range children {
		headerMap[key] = []string{child.Data().(string)}
	}
	return headerMap
}

func (ep *EndpointConfig) RequestMethod() string {
	method := strings.ToUpper(ep.GetString("method"))
	return method
}

func (ep *EndpointConfig) RequestURL() string {
	return ep.GetString("url")
}
