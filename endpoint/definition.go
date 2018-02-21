package endpoint

import (
	"github.com/5sigma/spyder/faker"
	"github.com/Jeffail/gabs"
)

// HeaderDefinition contains meta data about the request headers.
type HeaderDefinition struct {
	Key         string
	Example     string
	Description string
}

// ExampleResponse - Generates an example JSON response from the definition
// provided in the endpoint.
func (ep *EndpointConfig) ExampleResponse() string {
	definition := ep.ResponseDefinition()
	example := buildNode(definition, true).(*gabs.Container)
	return example.StringIndent("", "  ")
}

// ExampleRequest - Generates an example JSON request from the definition
// provided in the endpoint.
func (ep *EndpointConfig) ExampleRequest() string {
	definition := ep.RequestDefinition()
	example := buildNode(definition, true).(*gabs.Container)
	return example.StringIndent("", "  ")
}

func buildNode(defNode *gabs.Container, root bool) interface{} {
	nodeType, _ := defNode.Path("type").Data().(string)
	switch nodeType {
	case "object":
		newNode := gabs.New()
		propNode, _ := defNode.S("properties").ChildrenMap()
		for key, childNode := range propNode {
			newNode.Set(buildNode(childNode, false), key)
		}
		if root {
			return newNode
		} else {
			return newNode.Data()
		}
	case "string":
		return getExample(defNode, "123")
	case "number":
		return getExample(defNode, 123)
	case "bool":
		return getExample(defNode, true)
	case "array":
		return []interface{}{buildNode(defNode.Path("items"), false)}
	}
	return gabs.New()
}

func getExample(node *gabs.Container, defVal interface{}) interface{} {
	if node.Exists("example") {
		if val, ok := node.Path("example").Data().(string); ok {
			return faker.ExpandFakes(val)
		} else {
			return val
		}
	} else {
		return defVal
	}
}

// HeaderDefinitions - returns a list of information about the headers required
// for the request.
func (ep *EndpointConfig) HeaderDefinitions() []HeaderDefinition {
	defs := []HeaderDefinition{}
	childMap, _ := ep.json.Path("definition.headers").ChildrenMap()
	for k, v := range childMap {
		desc, _ := v.S("description").Data().(string)
		example, _ := v.S("example").Data().(string)
		defs = append(defs, HeaderDefinition{
			Key:         k,
			Example:     example,
			Description: desc,
		})
	}
	return defs
}
