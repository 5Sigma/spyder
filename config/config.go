package config

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Config - The config object is used to store the application configuration.
// The package attempts instantiate two configs one for the local config and one
// for the global config.
type Config struct {
	Filename string
	Filepath string
	json     *gabs.Container
}

// LocalConfig - The local configuration read from config.local.json
var LocalConfig = loadConfigFile("spyder.local.json")

// GlobalConfig - The global configuration read from config.json
var GlobalConfig = loadConfigFile("spyder.json")

// TempConfig - A temporary config object that persists only for the current
// session.
var TempConfig = loadDefaultConfig()

// The path to the project root
var ProjectPath = getProjectPath()

// InMemory - When true the config will not write to the disk. This is used for
// testing.
var InMemory = false

// GetProjectPath - returns the project root path it looks for a spyder.config
// file in the current folder. If found the directory tree is walked up until it
// finds one. If one is still not found the current path is used. Once a
// configuration file is found the path can be overriden using the projectPath
// setting.
func getProjectPath() string {
	p := GetSetting("projectPath")
	if p != "" {
		return path.Join(GlobalConfig.Filepath, p)
	} else {
		return "."
	}
}

func getConfigPath(p string) string {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return p
	} else {
		newPath, err := filepath.Abs(path.Join(filepath.Dir(p), "..", filepath.Base(p)))
		if err != nil {
			return ""
		}
		if newPath == p {
			return "."
		}
		if stat, err := os.Stat(filepath.Dir(newPath)); err == nil && stat.IsDir() {
			return getConfigPath(newPath)
		} else {
			return ""
		}
	}
	return ""
}

// VariableExists - Checks if a variable exists in either config.
func VariableExists(str string) bool {
	if TempConfig.VariableExists(str) {
		return true
	}
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
	if TempConfig.VariableExists(str) {
		return TempConfig.GetVariable(str)
	}
	if LocalConfig.VariableExists(str) {
		return LocalConfig.GetVariable(str)
	}
	if GlobalConfig.VariableExists(str) {
		return GlobalConfig.GetVariable(str)
	}
	return ""
}

// GetSetting - Returns the value of a variable from either config. Priority
// goes to the local config.
func GetSetting(str string) string {
	if LocalConfig.SettingExists(str) {
		return LocalConfig.GetSetting(str)
	}
	if GlobalConfig.SettingExists(str) {
		return GlobalConfig.GetSetting(str)
	}
	return ""
}

// GetSetting - Returns the value of a variable from either config. Priority
// goes to the local config. If no value is found the default is returned.
func GetSettingDefault(str, def string) string {
	res := GetSetting(str)
	if res == "" {
		return def
	}
	return res
}

// ExpandString - Given a string with a variable inside it. The string will be
// expanded and the variable placeholders replaced with variables from either
// config. Priority goes to the local config.
func ExpandString(str string) string {
	str = TempConfig.ExpandString(str)
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

	configPath := getConfigPath(filename)

	if configPath != "" {
		bytes, _ := ioutil.ReadFile(configPath)
		if strings.TrimSpace(string(bytes)) == "" {
			c = loadDefaultConfig()
			c.Filename = filename
			return c
		}
		c = LoadConfig(bytes)
		c.Filename = filename
		c.Filepath = filepath.Dir(configPath)
		return c
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

// SettingExists - returns true if a setting is specified in the config.
func (c *Config) SettingExists(name string) bool {
	return c.json.ExistsP(name)
}

// GetSetting - Returns a setting as a string from the path specified in the
// config.
func (c *Config) GetSetting(name string) string {
	v, _ := c.json.Path(name).Data().(string)
	return v
}

// GetSetting - Returns a setting as a string from the path specified in the
// config. If the setting does not exist the default value is returned.
func (c *Config) GetSettingDefault(name, def string) string {
	if c.SettingExists(name) {
		v, _ := c.json.Path(name).Data().(string)
		return v
	} else {
		return def
	}
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
