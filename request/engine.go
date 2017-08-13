package request

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	_ "github.com/robertkrimen/otto/underscore"
	"io/ioutil"
	"net/http"
)

// ScriptEngine - A scripting engine used to execute hook scripts during the
// request process. It is built to execute arbitrary Javascript.
type ScriptEngine struct {
	VM             *otto.Otto
	Constants      map[string]string
	AssetPath      string
	Response       *Response
	EndpointConfig *endpoint.EndpointConfig
	Payload        []byte
	Request        *http.Request
	Debug          string
}

// NewScriptEngine - Generates a new script engine.
func NewScriptEngine(endpointConfig *endpoint.EndpointConfig) *ScriptEngine {
	vm := otto.New()

	eng := &ScriptEngine{
		VM:             vm,
		EndpointConfig: endpointConfig,
	}

	varObj, _ := vm.Object("$variables = {}")
	varObj.Set("set", eng.setLocalVar)
	varObj.Set("get", eng.getVar)

	vm.Set("$debug", eng.setDebug)
	vm.Set("$hmac", eng.hmac)

	reqObj, _ := eng.VM.Object("$request = {}")
	requestBytes := endpointConfig.RequestData()
	reqObj.Set("body", string(requestBytes))
	reqObj.Set("contentLength", len(requestBytes))
	headersObj, _ := eng.VM.Object(`$request.headers = {}`)
	headersObj.Set("get", eng.getReqHeader)
	headersObj.Set("set", eng.setReqHeader)
	reqObj.Set("setBody", eng.setPayload)

	return eng
}

//ExecuteFile - Executes a script contianed in a file.
func (eng *ScriptEngine) ExecuteFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = eng.Execute(string(data))
	if err != nil {
		return err
	}
	return nil
}

// SetResponse - Used to set the web respone on the engine. It also makes this
// available to the script.
func (eng *ScriptEngine) SetResponse(res *Response) {
	eng.Response = res
	requestObj, _ := eng.VM.Object(`$request = {}`)
	requestObj.Set("url", res.Request.URL.String())
	requestObj.Set("contentLength", res.Request.ContentLength)

	responseObj, _ := eng.VM.Object(`$response = {}`)
	responseObj.Set("contentLength", res.Response.ContentLength)
	responseObj.Set("body", string(res.Content))
	eng.VM.Object(`$response.headers = {}`)
	responseObj.Set("get", eng.getResHeader)
}

// SetRequest - Sets the request on the engine. This also builds the functions
// to expose the request to scripts.
func (engine *ScriptEngine) SetRequest(request *http.Request) {
	engine.Request = request
	reqVal, _ := engine.VM.Get("$request")
	reqObj := reqVal.Object()
	reqObj.Set("contentLength", request.ContentLength)
}

// SetRequestConfig - Sets various properties of the exposed $request object
// using the endpoint configuration.
func (eng *ScriptEngine) SetEndpointConfig(epConf *endpoint.EndpointConfig) {
}

// SetPayload - Sets the request payload on the engine. Also exposes it to
// scripts within the request object.
func (engine *ScriptEngine) SetPayload(payload []byte) {
	engine.Payload = payload
	reqVal, _ := engine.VM.Get("$request")
	reqObj := reqVal.Object()
	reqObj.Set("body", string(payload))
}

//Execute - Executes a Javascript.
func (eng *ScriptEngine) Execute(script string) error {
	_, err := eng.VM.Run(script)
	if err != nil {
		println(err.Error())
	}
	return err
}

// Validate - Validates that the Javascript is valid.
func (eng *ScriptEngine) Validate(script string) error {
	_, err := parser.ParseFile(nil, "", script, 0)
	return err
}

// jsThrow - Used to throw javascript errors from Go.
func jsThrow(call otto.FunctionCall, err error) {
	call.Otto.Call("new Error", nil, err.Error())
}

// setLocalVar - Sets a config value in the local config.
func (engine *ScriptEngine) setLocalVar(call otto.FunctionCall) otto.Value {
	key, _ := call.Argument(0).ToString()
	value, _ := call.Argument(1).ToString()
	config.LocalConfig.SetVariable(key, value)
	return otto.Value{}
}

// getVar - Returns a config variable.
func (engine *ScriptEngine) getVar(call otto.FunctionCall) otto.Value {
	key, _ := call.Argument(0).ToString()
	if config.LocalConfig.VariableExists(key) {
		v := config.LocalConfig.GetVariable(key)
		ov, _ := otto.ToValue(v)
		return ov
	}
	if config.GlobalConfig.VariableExists(key) {
		v := config.LocalConfig.GetVariable(key)
		ov, _ := otto.ToValue(v)
		return ov
	}
	return otto.Value{}
}

// getPayload - Returns the request payload.
func (engine *ScriptEngine) getPayload(call otto.FunctionCall) otto.Value {
	ov, _ := otto.ToValue(string(engine.Payload))
	return ov
}

// setPayload - Sets the payload value on the engine.
func (engine *ScriptEngine) setPayload(call otto.FunctionCall) otto.Value {
	val, _ := call.Argument(0).ToString()
	engine.Payload = []byte(val)
	return otto.Value{}
}

// hmac - Calcualtes an HMAC signature using SHA256
func (engine *ScriptEngine) hmac(call otto.FunctionCall) otto.Value {
	secretStr, _ := call.Argument(0).ToString()
	payloadStr, _ := call.Argument(1).ToString()
	secret := []byte(secretStr)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(payloadStr))
	sigBytes := h.Sum(nil)
	v, _ := otto.ToValue(hex.EncodeToString(sigBytes))
	return v
}

// getReqHeader - Returns a header from the request.
func (engine *ScriptEngine) getReqHeader(call otto.FunctionCall) otto.Value {
	headerName, _ := call.Argument(0).ToString()
	val := engine.Request.Header.Get(headerName)
	v, _ := otto.ToValue(val)
	return v
}

// setReqHeader - sets a header on the request.
func (engine *ScriptEngine) setReqHeader(call otto.FunctionCall) otto.Value {
	headerName, _ := call.Argument(0).ToString()
	headerValue, _ := call.Argument(1).ToString()
	if engine.Request != nil {
		engine.Request.Header.Set(headerName, headerValue)
	} else {
		engine.EndpointConfig.Headers[headerName] = []string{headerValue}
	}
	return otto.Value{}
}

// getResHeader - Returns a header from the response.
func (engine *ScriptEngine) getResHeader(call otto.FunctionCall) otto.Value {
	headerName, _ := call.Argument(0).ToString()
	val := engine.Response.Response.Header.Get(headerName)
	v, _ := otto.ToValue(val)
	return v
}

// setDebug - Sets the debug value on the engine. Used for testing.
func (engine *ScriptEngine) setDebug(call otto.FunctionCall) otto.Value {
	val, _ := call.Argument(0).ToString()
	engine.Debug = val
	return otto.Value{}
}
