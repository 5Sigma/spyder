package request

import (
	"bytes"
	"fmt"
	"github.com/5sigma/spyder/endpoint"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
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
		requestData = config.GetJSONBytes("data")
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequest(config.Method,
		config.Url, bytes.NewReader(requestData))

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

func generateGetUrl(urlString string, postData map[string]string) string {
	baseUrl, err := url.Parse(urlString)
	if err != nil {
		return urlString
	}

	params := url.Values{}

	for k, v := range postData {
		if strings.HasPrefix(postData[k], "$") {
			varName := strings.Replace(postData[k], "$", "", 1)
			varValue := viper.GetString(fmt.Sprintf("variables.%s", varName))
			params.Add(k, varValue)
		} else {
			params.Add(k, v)
		}
	}

	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}
