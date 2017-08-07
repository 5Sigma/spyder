package config

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"os"
)

type (
	Config struct {
		filename string
		json     *gabs.Container
	}
)

var LocalConfig = loadConfig("config.local.json")
var GlobalConfig = loadConfig("config.json")

func loadConfig(filename string) *Config {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		bytes, _ := ioutil.ReadFile(filename)
		json, _ := gabs.ParseJSON(bytes)
		config := &Config{
			filename: filename,
			json:     json,
		}
		return config
	} else {
		return loadDefaultConfig(filename)
	}
}

func loadDefaultConfig(filename string) *Config {
	c := &Config{
		filename: filename,
		json:     gabs.New(),
	}

	c.json.Object("variables")
	return c
}

func (c *Config) Write() {
	ioutil.WriteFile(c.filename, c.json.BytesIndent("", "  "), os.ModePerm)
}

func (c *Config) GetVariable(path string) string {
	v, _ := c.json.Path(fmt.Sprintf("variables.%s", path)).Data().(string)
	return v
}

func (c *Config) SetVariable(path, value string) {
	c.json.SetP(value, fmt.Sprintf("variables.%s", path))
	c.Write()
}

func (c *Config) VariableExists(path string) bool {
	return c.json.ExistsP(fmt.Sprintf("variables.%s", path))
}
