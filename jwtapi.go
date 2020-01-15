package apiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var maxRetry = 2

type ijwtapi interface {
	GetTokenRequestData() TokenRequest
	GetTokenURI() string
	GetRefreshTokenURI() string
	DebugEnabled() bool
	InsecureSSLEnabled() bool
	GetTimeout() time.Duration
	GetToken() Token

	getClient() *http.Client
	logMsg(methodname, format string, msg ...interface{})
}

//JwtAPI provide functions to call JWT protected APIs by setting Access-Token in request Authorization header
type JwtAPI struct {
	TokenRequestData TokenRequest
	token            Token
	TokenURI         string
	RefreshTokenURI  string
	AllowInsecureSSL bool
	Timeout          time.Duration
	Debug            bool
}

//SJwtAPI allow to maek calls to JWT protected Structured APIs by setting Access-Token in request Authorization header.
//Called structured APIs are expected to return response as APIReuslt{}
type SJwtAPI struct {
	TokenRequestData TokenRequest
	token            Token
	TokenURI         string
	RefreshTokenURI  string
	AllowInsecureSSL bool
	Timeout          time.Duration
	Debug            bool
}

//Token is returned after successfull request to token or refreshtoken endpoints
type Token struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

//TokenRequest is used to pass credential to auth-server to get new token
type TokenRequest struct {
	ClientID     string
	ClientSecret string
	Scopes       string
	AppUserID    string
	RefreshToken string
}

//RefreshToken is used to get New AccessToken
type refreshToken struct {
	RefreshToken string
}

func (j JwtAPI) GetTokenRequestData() TokenRequest {
	return j.TokenRequestData
}
func (j JwtAPI) GetTokenURI() string {
	return j.TokenURI
}
func (j JwtAPI) GetRefreshTokenURI() string {
	return j.RefreshTokenURI
}
func (j JwtAPI) DebugEnabled() bool {
	return j.Debug
}
func (j JwtAPI) GetToken() Token {
	return j.token
}
func (j JwtAPI) InsecureSSLEnabled() bool {
	return j.AllowInsecureSSL
}
func (j JwtAPI) GetTimeout() time.Duration {
	return j.Timeout
}
func (j JwtAPI) getClient() *http.Client {
	return getClient(j.InsecureSSLEnabled(), j.GetTimeout())
}
func (j JwtAPI) logMsg(methodname, format string, msg ...interface{}) {
	logMsg(j.DebugEnabled(), methodname, format, msg...)
}

//
// ---------------------
//

func (sj SJwtAPI) GetTokenRequestData() TokenRequest {
	return sj.TokenRequestData
}
func (sj SJwtAPI) GetTokenURI() string {
	return sj.TokenURI
}
func (sj SJwtAPI) GetRefreshTokenURI() string {
	return sj.RefreshTokenURI
}
func (sj SJwtAPI) DebugEnabled() bool {
	return sj.Debug
}
func (sj SJwtAPI) GetToken() Token {
	return sj.token
}
func (sj SJwtAPI) InsecureSSLEnabled() bool {
	return sj.AllowInsecureSSL
}
func (sj SJwtAPI) GetTimeout() time.Duration {
	return sj.Timeout
}
func (sj SJwtAPI) getClient() *http.Client {
	return getClient(sj.InsecureSSLEnabled(), sj.GetTimeout())
}
func (sj SJwtAPI) logMsg(methodname, format string, msg ...interface{}) {
	logMsg(sj.DebugEnabled(), methodname, format, msg...)
}

func getClient(allowInsecureSSL bool, timeout time.Duration) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: allowInsecureSSL},
	}

	// default timeout (if not set by client)
	timeoutInSec := 10

	if timeout.Seconds() > 0.1 {
		// client set timeout, so use it
		timeoutInSec = int(timeout.Seconds())
	}

	client := &http.Client{
		Timeout:   time.Second * time.Duration(timeoutInSec),
		Transport: tr,
	}
	return client
}

func logMsg(debug bool, methodname, format string, msg ...interface{}) {
	if !debug {
		return
	}
	log.Printf("DEBUG: [%s] [%s]\n", methodname, fmt.Sprintf(format, msg...))
}

func requestTokenByLogin(j ijwtapi) (Token, error) {
	var token Token

	j.logMsg("RequestTokenByLogin", "%s", "Requesting new token through login")
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(j.GetTokenRequestData())

	r, err := http.NewRequest(http.MethodPost, j.GetTokenURI(), b)
	if err != nil {
		return token, err
	}
	r.Header.Set("Content-Type", "application/json")

	client := j.getClient()
	resp, err := client.Do(r)
	if err != nil {
		if resp != nil {
			j.logMsg("RequestTokenByLogin", "failed to get new token by login (%d): %v", resp.StatusCode, err)
			resp.Body.Close()
			return token, fmt.Errorf("failed to get new token by login (%d): %v", resp.StatusCode, err)
		}
		j.logMsg("RequestTokenByLogin", "failed to get new token by login: %v\n", err)
		return token, fmt.Errorf("failed to get new token by login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseData, err := ioutil.ReadAll(resp.Body)
		respStr := ""
		if err == nil {
			respStr = string(responseData)
		}
		j.logMsg("RequestTokenByLogin", "failed to get new token[2] by login (%d): %v\n", resp.StatusCode, err)
		return token, fmt.Errorf("failed to get new token by login (%d): %s", resp.StatusCode, respStr)
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return token, fmt.Errorf("failed to extract new token: %v", err)
	}

	return token, nil
}

func requestTokenByRefreshToken(j ijwtapi, rtoken string) (Token, error) {
	var token Token

	j.logMsg("RequestTokenByRefreshToken", "Requesting new token through refresh-token (%v)", rtoken)

	u := refreshToken{RefreshToken: rtoken}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	r, err := http.NewRequest(http.MethodPost, j.GetRefreshTokenURI(), b)
	if err != nil {
		return token, err
	}
	r.Header.Set("Content-Type", "application/json")

	client := j.getClient()
	resp, err := client.Do(r)
	if err != nil {
		if resp != nil {
			j.logMsg("RequestTokenByRefreshToken", "failed to get new token by refreshtoken (%d): %v\n", resp.StatusCode, err)
			resp.Body.Close()
			return token, fmt.Errorf("failed to get new token by refreshtoken (%d): %v", resp.StatusCode, err)
		}
		j.logMsg("RequestTokenByRefreshToken", "failed to get new token by refreshtoken: %v\n", err)
		return token, fmt.Errorf("failed to get new token by refreshtoken: %v", err)
	}
	defer resp.Body.Close()

	// read response string

	if resp.StatusCode != 200 {
		// //DEBUG
		// responseData, err := ioutil.ReadAll(r.Body)
		// defer r.Body.Close()
		// if err == nil {
		// 	log.Printf("\tRefreshToke Response: %s\n", responseData)
		// } else {
		// 	log.Println("\tRefreshToke Response: not-readable")
		// }
		// //DEBUG end

		//return "", fmt.Errorf("failed to get new token by refresh-token (%d)", r.StatusCode)

		// Possibly refresh-token expired or there is scope mismatch
		//   Try to get a fresh AccessToken by login
		return requestTokenByLogin(j)
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		j.logMsg("RequestTokenByRefreshToken", "Token unmarshal error: %v", err)
		return token, err
	}

	return token, nil
}

//makeRequest makes http request for given url with given method
func makeRequest(j ijwtapi, method, apiurl string, body io.Reader) (*http.Response, error) {
	retry := 0
	connFailRetry := 0
	//Create []byte buffer from body - so it can be passed in further retries
	var buf []byte
	if body != nil {
		buf, _ = ioutil.ReadAll(body)
	}

callapi:
	j.logMsg("makeRequest", "Retry[%d], API: %s\n\tBody: %s", retry, apiurl, buf)

	r, err := http.NewRequest(method, apiurl, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "bearer "+j.GetToken().AccessToken)
	r.Header.Set("Content-Type", "application/json")

	//client := &http.Client{}
	client := j.getClient()
	resp, err := client.Do(r)
	if err != nil {
		errmsg := err.Error()
		j.logMsg("makeRequest", "Api-Error: %s", errmsg)
		if (strings.Contains(errmsg, "getsockopt: connection")) && (connFailRetry < maxRetry) {
			// couldn't connect to remote API server, connection failed, try again
			time.Sleep(time.Millisecond * 500)
			connFailRetry++
			if resp != nil {
				resp.Body.Close()
			}
			goto callapi
		}
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// //DEBUG
		// if retry > 0 {
		// 	tmpres := handleNotOK(resp)
		// 	log.Printf("\tRetried API, got error: %s\n", tmpres.ErrText)
		// }
		// //DEBUG end

		if (retry < maxRetry) && (resp.StatusCode == http.StatusUnauthorized) {
			resp.Body.Close()

			retry++
			requestTokenByRefreshToken(j, j.GetToken().RefreshToken)
			// again try to call same API
			goto callapi
		}
	}
	return resp, nil
}