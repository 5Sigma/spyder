package request

import (
	"bytes"
	"github.com/5sigma/spyder/endpoint"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

type (
	Response struct {
		Response    *http.Response
		RequestTime time.Duration
		Content     []byte
		Request     *http.Request
		Payload     []byte
	}
)

func Do(config *endpoint.EndpointConfig) (*Response, error) {
	var (
		client       = &http.Client{}
		err          error
		req          *http.Request
		requestData  = []byte{}
		scriptEngine = NewScriptEngine(config)
	)

	if config.Method == "POST" || config.Method == "HEAD" {
		requestData = config.RequestData()
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequest(config.Method,
		config.RequestURL(), bytes.NewReader(requestData))

	req.Header = config.Headers()

	scriptEngine.Request = req

	// Perform transformation
	// executes scripts with the payload set so that they are capable of modifying
	// it before the request is made.
	if len(config.Transform) > 0 {
		for _, transform := range config.Transform {
			scriptPath := path.Join("scripts", transform) + ".js"
			requestData = scriptEngine.ExecuteTransform(scriptPath, requestData)
		}
	}

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
	if len(config.OnComplete) > 0 {
		scriptEngine.SetResponse(res)
		for _, onComplete := range config.OnComplete {
			scriptPath := path.Join("scripts", onComplete) + ".js"
			err := scriptEngine.ExecuteFile(scriptPath)
			if err != nil {
				return res, err
			}
		}
	}

	return res, nil
}

func (res *Response) IsResponseJSON() bool {
	contentType := strings.ToLower(res.Response.Header.Get("Content-Type"))
	return strings.Contains(contentType, "json")
}
