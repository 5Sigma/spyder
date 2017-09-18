package docgen

import (
	"bytes"
	"github.com/5sigma/spyder/endpoint"
	"text/template"
)

type EndpointData struct {
	Name            string
	Method          string
	ResponseExample string
}

var epTemplate = `
	{{.Name}} - {{.Method}}
	Response:

	{{.ResponseExample}}
`

func EndpointSection(path string, ep *endpoint.EndpointConfig) (string, error) {
	data := EndpointData{
		Name:            path,
		Method:          ep.Method,
		ResponseExample: ep.ExampleResponse(),
	}
	tmpl, err := template.New("endpoint").Parse(epTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
