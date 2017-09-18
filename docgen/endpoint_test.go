package docgen

import (
	"github.com/5sigma/spyder/endpoint"
	"testing"
)

func ep() *endpoint.EndpointConfig {
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
	ep, err := endpoint.LoadBytes([]byte(config))
	if err != nil {
		panic(err)
	}
	return ep
}

func TestEndpointSection(t *testing.T) {
	res, _ := EndpointSection("testing/test", ep())
	t.Error(res)
}
