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

	"github.com/chrikoch/go-fritzbox/config"
)

type fritzBoxSessionInfo struct {
	XMLName   xml.Name `xml:"SessionInfo"`
	SID       string
	Challenge string `xml:"Challenge"`
	BlockTime int64
}

// Authenticator is used to retrive a valid session id
type Authenticator struct {
	Config                config.Config
	sessionID             string
	sessionIDCreationTime time.Time
}

const invalidSessionID = "0000000000000000"
const sessionLifeTime = time.Minute * 60

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
		}
		return r

	}, s)

}

// SessionID returns a valid sessionID, valid for at least minRemainingLifeTime
func (a *Authenticator) SessionID(minRemainingLifeTime time.Duration) string {

	//sessionId set and EndOfLife is still longer than minimum remainder?
	if a.sessionID != invalidSessionID &&
		(a.sessionIDCreationTime.Add(sessionLifeTime).After(time.Now().Add(minRemainingLifeTime))) {
		return a.sessionID
	}

	//we have to get a new sessionID!
	id, err := a.newSessionID()
	if err != nil {
		log.Println(err)
		a.sessionID = invalidSessionID
		a.sessionIDCreationTime = time.Time{}

	} else {
		a.sessionID = id
		a.sessionIDCreationTime = time.Now()
	}

	return a.sessionID
}

func (a *Authenticator) newSessionID() (string, error) {
	req, err := http.NewRequest("GET", a.Config.LoginURL(), nil)
	if err != nil {
		log.Println(err)
		return invalidSessionID, err
	}

	query := req.URL.Query()

	response, err := a.response()
	if err != nil {
		log.Println(err)
		return invalidSessionID, err
	}

	query.Add("response", response)
	query.Add("username", "")

	req.URL.RawQuery = query.Encode()

	var client http.Client

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return invalidSessionID, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return invalidSessionID, err
	}

	var info fritzBoxSessionInfo

	err = xml.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return invalidSessionID, err
	}

	return info.SID, nil

}

func (a *Authenticator) newChallenge() (challenge string, err error) {
	resp, err := http.Get(a.Config.LoginURL())
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var info fritzBoxSessionInfo
	err = xml.Unmarshal(body, &info)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return info.Challenge, nil
}
