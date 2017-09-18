package endpoint

import (
	"github.com/xeipuuv/gojsonschema"
)

func (ep *EndpointConfig) ValidateResponse(response string) (*gojsonschema.Result, error) {
	var (
		schema string
	)

	schema = ep.ResponseDefinition().String()

	schemaLoader := gojsonschema.NewStringLoader(schema)
	requestLoader := gojsonschema.NewStringLoader(response)

	return gojsonschema.Validate(schemaLoader, requestLoader)
}
