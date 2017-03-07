package home_automation

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chrikoch/go-fritzbox/config"
)

type HomeAutomation struct {
	Config config.Config
}

//Returns current powerconsumption auf AIN as mW
func (h *HomeAutomation) CurrentPowerConsumption(sessionId, ain string) (power int, err error) {
	req, err := http.NewRequest("GET", h.Config.HomeAutomationUrl(), nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	query := req.URL.Query()

	query.Add("switchcmd", "getswitchpower")
	query.Add("sid", sessionId)
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

func (h *HomeAutomation) SwitchList(sessionId string) {
	req, err := http.NewRequest("GET", h.Config.HomeAutomationUrl(), nil)
	if err != nil {
		log.Println(err)
		return //InvalidSessionId, err
	}

	query := req.URL.Query()

	query.Add("switchcmd", "getswitchlist")
	query.Add("sid", sessionId)

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
