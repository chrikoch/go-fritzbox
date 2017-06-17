package home_automation

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chrikoch/go-fritzbox/config"
)

//HomeAutomation is used to contact the home automation webservices of FritzBox
type HomeAutomation struct {
	Config config.Config //must be a valid config
}

//CurrentPowerConsumption returns current powerconsumption auf AIN as mW
func (h *HomeAutomation) CurrentPowerConsumption(sessionID, ain string) (power int, err error) {
	req, err := http.NewRequest("GET", h.Config.HomeAutomationURL(), nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	query := req.URL.Query()

	query.Add("switchcmd", "getswitchpower")
	query.Add("sid", sessionID)
	query.Add("ain", ain)

	req.URL.RawQuery = query.Encode()

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(body)))
}

//SwitchList logs a list of available switches
func (h *HomeAutomation) SwitchList(sessionID string) {
	req, err := http.NewRequest("GET", h.Config.HomeAutomationURL(), nil)
	if err != nil {
		log.Println(err)
		return //InvalidSessionId, err
	}

	query := req.URL.Query()

	query.Add("switchcmd", "getswitchlist")
	query.Add("sid", sessionID)

	req.URL.RawQuery = query.Encode()

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return //InvalidSessionId, err
	}

	defer resp.Body.Close()

	log.Printf("Got HTTP result %v.", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return //InvalidSessionId, err
	}

	log.Println(string(body))

}
