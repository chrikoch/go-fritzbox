package auth

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
	Config                config.Config
	sessionId             string
	sessionIdCreationTime time.Time
}

const InvalidSessionId = "0000000000000000"
const SessionLifeTime = time.Minute * 60

func (a *Authenticator) response() (response string, err error) {
	challenge, err := a.newChallenge()

	if err != nil {
		return "", err
	}

	return calculateResponse(challenge, a.Config.Password), nil
}

func calculateResponse(challenge, password string) (response string) {
	unhashed := challenge + "-" + password

	//now replace every rune > 255
	unhashed = replaceInvalidChallengeRunes(unhashed)

	unhashedUtf16 := utf8stringToUtf16Le(unhashed)
	hashed := md5.Sum(unhashedUtf16)

	response = fmt.Sprintf("%x", hashed)

	return challenge + "-" + response
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

func (a *Authenticator) SessionId(minRemainingLifeTime time.Duration) string {

	//sessionId set and EndOfLife is still longer than minimum remainder?
	if a.sessionId != InvalidSessionId &&
		(a.sessionIdCreationTime.Add(SessionLifeTime).After(time.Now().Add(minRemainingLifeTime))) {
		return a.sessionId
	}

	//we have to get a new sessionID!
	id, err := a.newSessionId()
	if err != nil {
		log.Println(err)
		a.sessionId = InvalidSessionId
		a.sessionIdCreationTime = time.Time{}

	} else {
		a.sessionId = id
		a.sessionIdCreationTime = time.Now()
	}

	return a.sessionId
}

func (a *Authenticator) newSessionId() (string, error) {
	req, err := http.NewRequest("GET", a.Config.LoginUrl(), nil)
	if err != nil {
		log.Println(err)
		return InvalidSessionId, err
	}

	query := req.URL.Query()

	response, err := a.response()
	if err != nil {
		log.Println(err)
		return InvalidSessionId, err
	}

	query.Add("response", response)
	query.Add("username", "")

	req.URL.RawQuery = query.Encode()

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return InvalidSessionId, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return InvalidSessionId, err
	}

	var info FritzBoxSessionInfo

	err = xml.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return InvalidSessionId, err
	}

	return info.SID, nil

}

func (a *Authenticator) newChallenge() (challenge string, err error) {
	resp, err := http.Get(a.Config.LoginUrl())
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var info FritzBoxSessionInfo
	err = xml.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return info.Challenge, nil
}
