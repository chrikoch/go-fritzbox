package auth

import (
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/chri.koch/fritzbox_util/config"
)

type Authenticator struct {
	Config config.Config
}

func (a *Authenticator) GetChallenge() {
	resp, err := http.Get("http://" + a.Config.BoxURL + "/login_sid.lua")
	if err != nil {
		// handle error
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("%v\n", string(body))
}
