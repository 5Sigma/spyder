package config

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Config - The config object is used to store the application configuration.
// The package attempts instantiate two configs one for the local config and one
// for the global config.
type Config struct {
	Filename string
	json     *gabs.Container
}

// LocalConfig - The local configuration read from config.local.json
var LocalConfig = loadConfigFile("config.local.json")

// GlobalConfig - The global configuration read from config.json
var GlobalConfig = loadConfigFile("config.json")

// The path to the project root
var ProjectPath = "."

// InMemory - When true the config will not write to the disk. This is used for
// testing.
var InMemory = false

// VariableExists - Checks if a variable exists in either config.
func VariableExists(str string) bool {
	if LocalConfig.VariableExists(str) {
		return true
	}
	if GlobalConfig.VariableExists(str) {
		return true
	}
	return false
}

// GetVariable - Returns the value of a variable from either config. Priority
// goes to the local config.
func GetVariable(str string) string {
	if LocalConfig.VariableExists(str) {
		return LocalConfig.GetVariable(str)
	}
	if GlobalConfig.VariableExists(str) {
		return GlobalConfig.GetVariable(str)
	}
	return ""
}

// ExpandString - Given a string with a variable inside it. The string will be
// expanded and the variable placeholders replaced with variables from either
// config. Priority goes to the local config.
func ExpandString(str string) string {
	str = LocalConfig.ExpandString(str)
	str = GlobalConfig.ExpandString(str)
	return str
}

// LoadConfigFile - Loads a config from a file on the disk.
func loadConfigFile(filename string) *Config {
	var (
		c *Config
	)
	if InMemory {
		return loadDefaultConfig()
	}

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		bytes, _ := ioutil.ReadFile(filename)
		if strings.TrimSpace(string(bytes)) == "" {
			c = loadDefaultConfig()
		}
		c = LoadConfig(bytes)
	}
	c = loadDefaultConfig()
	c.Filename = filename
	return c
}

// Loads a config from a byte array.
func LoadConfig(bytes []byte) *Config {
	json, err := gabs.ParseJSON(bytes)
	if err != nil {
		println("\n *** Could not parse config file. ***\n")
		println(err.Error())
		return loadDefaultConfig()
	}
	config := &Config{
		json: json,
	}
	return config
}

// LoadDefaultConfig - Loads a default empty configuration. This is used as a
// fallback when a config file does not exist or an error occures when reading
// it.
func loadDefaultConfig() *Config {
	c := &Config{
		json: gabs.New(),
	}
	c.json.Object("variables")
	return c
}

// Write - writes the config to the disk. Uses the path specified in the
// Filename property.
func (c *Config) Write() {
	if InMemory {
		return
	}
	ioutil.WriteFile(c.Filename, c.json.BytesIndent("", "  "), os.ModePerm)
}

// GetVariable returns a saved variable from the config.
func (c *Config) GetVariable(path string) string {
	v, _ := c.json.Path(fmt.Sprintf("variables.%s", path)).Data().(string)
	return v
}

// SetVariable - sets a variable in the config.
func (c *Config) SetVariable(path, value string) {
	c.json.SetP(value, fmt.Sprintf("variables.%s", path))
	c.Write()
}

// VariableExists - Checks to see if a variable exists at the path specified.
func (c *Config) VariableExists(path string) bool {
	return c.json.ExistsP(fmt.Sprintf("variables.%s", path))
}

// GetSetting - Returns a setting as a string from the path specified in the
// config.
func (c *Config) GetSetting(name string) string {
	v, _ := c.json.Path(name).Data().(string)
	return v
}

// String - returns the config JSON as a string.
func (c *Config) String() string {
	return c.json.String()
}

// ExpandString - Given a string with a variable inside it. The string will be
// expanded and the variable placeholders replaced with variables from the
// config.
func (c *Config) ExpandString(str string) string {
	re := regexp.MustCompile(`\$([A-Za-z0-9]+)`)
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		v := c.GetVariable(match[1])
		if v != "" {
			str = strings.Replace(str, fmt.Sprintf("$%s", match[1]), v, 1)
		}
	}
	return str
}
