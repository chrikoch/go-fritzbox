//holds the configuration
package config

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
