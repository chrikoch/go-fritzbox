//holds the configuration
package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	BoxURL   string
	Password string
}

func (c *Config) LoginUrl() string {
	return "http://" + c.BoxURL + "/login_sid.lua"
}

func (c *Config) HomeAutomationUrl() string {
	return "http://" + c.BoxURL + "/webservices/homeautoswitch.lua"
}

//Read a config in JSON format from filname and return as Config struct
func New(filename string) (Config, error) {
	var c Config
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return c, err
	}

	if len(c.BoxURL) == 0 {
		c.BoxURL = "fritz.box"
	}

	return c, nil
}
