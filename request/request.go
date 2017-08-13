package request

import (
	"bytes"
	"fmt"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

// Do - Performs the request for a given endpoint configuration. This will
// expand all variables, fake data, and execute any transforms and onComplete
// scripts specified.
func Do(epConfig *endpoint.EndpointConfig) (*Response, error) {
	var (
		client       = &http.Client{}
		err          error
		req          *http.Request
		requestData  = []byte{}
		scriptEngine = NewScriptEngine(epConfig)
	)

	if epConfig.RequestMethod() == "POST" || epConfig.RequestMethod() == "HEAD" {
		requestData = epConfig.RequestData()
		if err != nil {
			return nil, err
		}
	}

	scriptEngine.SetPayload(requestData)

	// Perform transformation
	// executes scripts with the payload set so that they are capable of modifying
	// it before the request is made.
	if len(epConfig.Transform) > 0 {
		for _, transform := range epConfig.Transform {
			scriptPath := path.Join(config.ProjectPath, "scripts", transform) + ".js"
			scriptEngine.ExecuteFile(scriptPath)
			requestData = scriptEngine.Payload
		}
	}

	req, err = http.NewRequest(epConfig.RequestMethod(),
		epConfig.RequestURL(), bytes.NewReader(requestData))

	req.Header = epConfig.Headers

	scriptEngine.Request = req

	// Make the request and calculate its flight time.
	start := time.Now()
	rawResponse, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, err
	}

	// Read response data
	contentBytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	res := &Response{
		Response:    rawResponse,
		RequestTime: elapsed,
		Content:     contentBytes,
		Request:     rawResponse.Request,
		Payload:     requestData,
	}

	// Execute any post request scripts that are listed in the config
	if len(epConfig.OnComplete) > 0 {
		scriptEngine.SetResponse(res)
		for _, onComplete := range epConfig.OnComplete {
			scriptPath := path.Join(config.ProjectPath, "scripts", onComplete) + ".js"
			err := scriptEngine.ExecuteFile(scriptPath)
			if err != nil {
				return res, fmt.Errorf("Error parsing script (%s): %s",
					scriptPath, err.Error())
			}
		}
	}

	return res, nil
}
