package home_automation

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chrikoch/go-fritzbox/config"
)

type HomeAutomation struct {
	Config config.Config
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
