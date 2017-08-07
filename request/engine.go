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

type (
	ScriptEngine struct {
		VM             *otto.Otto
		Constants      map[string]string
		AssetPath      string
		Response       *Response
		EndpointConfig *endpoint.EndpointConfig
		Payload        []byte
		Request        *http.Request
	}
)

func NewScriptEngine(endpointConfig *endpoint.EndpointConfig) *ScriptEngine {

	vm := otto.New()

	eng := &ScriptEngine{
		VM:             vm,
		EndpointConfig: endpointConfig,
	}

	varObj, _ := vm.Object("$variables = {}")
	varObj.Set("set", eng.engine_setLocalVar)
	varObj.Set("get", eng.engine_getVar)

	payloadObj, _ := vm.Object("$payload = {}")
	payloadObj.Set("get", eng.getPayload)
	payloadObj.Set("set", eng.setPayload)

	headerObj, _ := vm.Object("$headers = {}")
	headerObj.Set("get", eng.getHeader)
	headerObj.Set("set", eng.setHeader)

	vm.Set("$hmac", eng.hmac)

	return eng
}

//read javascript file and execute
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

func (eng *ScriptEngine) SetResponse(res *Response) {
	eng.Response = res
	requestObj, _ := eng.VM.Object(`$request = {}`)
	requestObj.Set("url", res.Request.URL.String())
	requestObj.Set("contentLength", res.Request.ContentLength)

	responseObj, _ := eng.VM.Object(`$response = {}`)
	responseObj.Set("contentLength", res.Response.ContentLength)
	responseObj.Set("body", string(res.Content))
}

//executes commands
func (eng *ScriptEngine) Execute(script string) error {
	_, err := eng.VM.Run(script)
	if err != nil {
		println(err.Error())
	}
	return err
}

func (eng *ScriptEngine) ExecuteTransform(script string, payload []byte) []byte {
	eng.Payload = payload
	eng.ExecuteFile(script)
	return eng.Payload
}

func (eng *ScriptEngine) Validate(script string) error {
	_, err := parser.ParseFile(nil, "", script, 0)
	return err
}

func jsThrow(call otto.FunctionCall, err error) {
	value, _ := call.Otto.Call("new Error", nil, err.Error())
	panic(value)
}

func (engine *ScriptEngine) engine_setLocalVar(call otto.FunctionCall) otto.Value {
	key, _ := call.Argument(0).ToString()
	value, _ := call.Argument(1).ToString()
	config.LocalConfig.SetVariable(key, value)
	return otto.Value{}
}

func (engine *ScriptEngine) engine_getVar(call otto.FunctionCall) otto.Value {
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

func (engine *ScriptEngine) getPayload(call otto.FunctionCall) otto.Value {
	ov, _ := otto.ToValue(string(engine.Payload))
	return ov
}

func (engine *ScriptEngine) setPayload(call otto.FunctionCall) otto.Value {
	val, _ := call.Argument(0).ToString()
	engine.Payload = []byte(val)
	return otto.Value{}
}

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

func (engine *ScriptEngine) getHeader(call otto.FunctionCall) otto.Value {
	headerName, _ := call.Argument(0).ToString()
	val := engine.Request.Header.Get(headerName)
	v, _ := otto.ToValue(val)
	return v
}

func (engine *ScriptEngine) setHeader(call otto.FunctionCall) otto.Value {
	headerName, _ := call.Argument(0).ToString()
	headerValue, _ := call.Argument(1).ToString()
	engine.Request.Header.Set(headerName, headerValue)
	return otto.Value{}
}
