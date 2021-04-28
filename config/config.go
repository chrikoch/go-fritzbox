//Package config holds the configuration
package config

//Config holds the Config of the FritzBox
type Config struct {
	BoxURL   string //URL of FritzBox
	Username string //Username used to login on FritzBox
	Password string //Password used to login on FritzBox
}

//LoginURL returns to URL for the Login webservice
func (c *Config) LoginURL() string {
	return "http://" + c.BoxURL + "/login_sid.lua"
}

//HomeAutomationURL returns the URL of the home automation webservice
func (c *Config) HomeAutomationURL() string {
	return "http://" + c.BoxURL + "/webservices/homeautoswitch.lua"
}
