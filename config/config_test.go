package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InMemory = true
	retCode := m.Run()
	os.Exit(retCode)
}

func TestLoadConfig(t *testing.T) {
	configJson := `
		{
			"variables": {
				"one": "1"
			}
		}
	`
	config := LoadConfig([]byte(configJson))
	val := config.GetVariable("one")

	if val != "1" {
		t.Errorf("Variable not set correctly: %s", val)
	}

	if config.VariableExists("one") != true {
		t.Errorf("Variable existance not detected")
	}

	if config.VariableExists("two") == true {
		t.Errorf("Variable existance false positive")
	}
}

func TestSetVariable(t *testing.T) {
	config := loadDefaultConfig()
	config.SetVariable("one", "1")
	v := config.GetVariable("one")
	if v != "1" {
		t.Errorf("Variable not set correctly: %s", v)
	}
}

func TestGetVariable(t *testing.T) {
	LocalConfig.SetVariable("getVal", "local")
	GlobalConfig.SetVariable("getVal", "global")
	GlobalConfig.SetVariable("another", "global")
	str := ExpandString("Value is '$getVal'")
	expected := "Value is 'local'"
	if str != expected {
		t.Errorf("Expanded value not correct\nReceived: %s\nExepcted:%s", str,
			expected)
	}
}

func TestVariableExists(t *testing.T) {
	LocalConfig.SetVariable("getVal", "local")
	GlobalConfig.SetVariable("getVal", "global")
	GlobalConfig.SetVariable("another", "global")
	if !VariableExists("getVal") {
		t.Errorf("Variablle getVar should exist")
	}
	if !VariableExists("another") {
		t.Errorf("Variablle another should exist")
	}
	if VariableExists("nope") {
		t.Errorf("Variablle nope should not exist")
	}
}
