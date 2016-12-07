package auth

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"gitlab.com/chri.koch/fritzbox_util/config"
)

type FritzBoxSessionInfo struct {
	XMLName   xml.Name `xml:SessionInfo`
	SID       string
	Challenge string `xml:Challenge`
	BlockTime int64
}

type Authenticator struct {
	Config config.Config
}

func (a *Authenticator) GetResponse() (resonse string, err error) {
	challenge, err := a.GetNewChallenge()

	if err != nil {
		return "", nil
	}

	unhashed := challenge + "-" + a.Config.Password
	//now replace every rune > 255
	unhashed = replaceInvalidChallengeRunes(unhashed)

	log.Println(unhashed)
	return unhashed, nil
}

func utf8stringToUtf16Le(s string) []byte {
	var buffer bytes.Buffer
	utf16encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	writer := transform.NewWriter(&buffer, utf16encoder.NewEncoder())

	io.WriteString(writer, s)
	writer.Close()

	return buffer.Bytes()
}

func replaceInvalidChallengeRunes(s string) string {
	return strings.Map(func(r rune) rune {
		if r > 255 {
			return '.'
		} else {
			return r
		}
	}, s)

}

func (a *Authenticator) GetNewChallenge() (challenge string, err error) {
	resp, err := http.Get("http://" + a.Config.BoxURL + "/login_sid.lua")
	if err != nil {
		// handle error
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("%v\n", string(body))

	var info FritzBoxSessionInfo

	err = xml.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(info)

	log.Println(info.Challenge)

	return info.Challenge, nil
}
