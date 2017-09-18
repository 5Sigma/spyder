package endpoint

import (
	"github.com/Jeffail/gabs"
)

func (ep *EndpointConfig) ExampleResponse() string {
	definition := ep.ResponseDefinition()
	example := buildNode(definition).(*gabs.Container)
	return example.StringIndent("", "  ")
}

func buildNode(defNode *gabs.Container) interface{} {
	nodeType, _ := defNode.Path("type").Data().(string)
	switch nodeType {
	case "object":
		newNode := gabs.New()
		propNode := defNode.S("properties")
		setProperties(newNode, propNode)
		return newNode
	case "string":
		return "123"
	case "number":
		return 123
	}
	return gabs.New()
}

func setProperties(exNode, defNode *gabs.Container) {
	children, _ := defNode.ChildrenMap()
	for key, childNode := range children {
		nodeType, _ := childNode.Path("type").Data().(string)
		switch nodeType {
		case "string":
			exNode.Set("abc", key)
		case "number":
			exNode.Set(123, key)
		case "array":
			exNode.Array(key)
			if obj, ok := buildNode(childNode.Path("items")).(*gabs.Container); ok {
				exNode.ArrayAppend(obj.Data(), key)
			} else {
				exNode.ArrayAppend(buildNode(childNode.Path("items")), key)
			}
		}
	}
}
