package request

import (
	"github.com/5sigma/spyder/config"
	"github.com/5sigma/spyder/endpoint"
	"net/http"
	"testing"
)

func TestSetVariable(t *testing.T) {
	script := ` $variables.set('key', 'value'); `
	engine := NewScriptEngine(endpoint.New())
	engine.Execute(script)
	if config.LocalConfig.GetVariable("key") != "value" {
		t.Errorf("Config value not set: %s", config.LocalConfig.GetVariable("key"))
	}
}

func TestDebug(t *testing.T) {
	script := `$debug('debug'); `
	engine := NewScriptEngine(endpoint.New())
	engine.Execute(script)
	if engine.Debug != "debug" {
		t.Errorf("Debug not set: %s", engine.Debug)
	}
}

func TestGetVariable(t *testing.T) {
	config.LocalConfig.SetVariable("key", "test1")
	script := `$debug($variables.get('key')); `
	engine := NewScriptEngine(endpoint.New())
	engine.Execute(script)
	if engine.Debug != "test1" {
		t.Errorf("Config value not set: %s", engine.Debug)
	}
}

func TestPaylaod(t *testing.T) {
	script := `
		$request.setBody($request.body + ' world');
	`
	engine := NewScriptEngine(endpoint.New())
	engine.SetPayload([]byte(`hello`))
	engine.Execute(script)
	if string(engine.Payload) != "hello world" {
		t.Errorf("Payload not set: %s", engine.Payload)
	}
}

func TestHMAC(t *testing.T) {
	script := `$debug($hmac('secret', 'hello'))`
	engine := NewScriptEngine(endpoint.New())
	engine.Execute(script)
	expected := "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b"
	if engine.Debug != expected {
		t.Errorf("Hash missmatch:\nExpected: %s\nReceieved: %s",
			expected, engine.Debug)
	}
}

func TestHeaders(t *testing.T) {
	script := `
		header1 = $request.headers.get('test-header1');
		$request.headers.set('test-header2', header1);
	`
	engine := NewScriptEngine(endpoint.New())
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	req.Header.Set("test-header1", "myval")
	engine.Request = req
	engine.Execute(script)
	if req.Header.Get("test-header2") != "myval" {
		headerArr := req.Header.Get("test-header2")
		if len(headerArr) == 0 {
			t.Fatal("Header not present")
		}
		t.Errorf("Header not set: %s", req.Header.Get("test-header2")[0])
	}
}
