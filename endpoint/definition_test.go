package endpoint

import (
	"testing"
)

func endpoint() *EndpointConfig {
	config := `
	{
		"url": "http://www.test.com",
		"method": "GET",
		"definition": {
			"response": {
				"type": "object",
				"properties": {
					"errorMessage": {
						"type": "string"
					},
					"errorNumber": {
						"type": "number"
					},
					"data": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"name": {
									"type": "string"
								},
								"id": {
									"type": "number"
								},
								"stringArray": {
									"type": "array",
									"items": {
										"type": "string"
									}
								},
								"numberArray": {
									"type": "array",
									"items": {
										"type": "number"
									}
								}
							}
						}
					}
				}
			}
		}
	}
	`
	ep, err := LoadBytes([]byte(config))
	if err != nil {
		panic(err)
	}
	return ep
}

func TestExampleResponse(t *testing.T) {
	ep := endpoint()
	res := ep.ExampleResponse()
	// t.Error(res)
}
