package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// SetLogWriter Sets io.writer for logging to file
func (j *JwtAPI) SetLogWriter(w io.Writer) {
	multi := io.MultiWriter(w, os.Stdout)
	l := log.New(multi, "", log.LstdFlags)
	l.Println("Log output set to StdOut and writer both")
	j.logger = l
}

// RequestTokenByCred call Token endpoint to get new token by passing TokenRequest data.
// It set token to JwtAPI instance for subsequent calls through same instance.
func (j *JwtAPI) RequestTokenByCred() (Token, error) {
	err := j.requestTokenByLogin()
	if err != nil {
		return Token{}, err
	}
	//j.logDebug("RequestTokenByRefreshToken", "API Object's refresh-token is (%s)", j.GetToken().RefreshToken)
	return j.GetToken(), nil
}

// RequestTokenByRefreshToken call Token endpoint to get new token by passing existing refresh-token.
// It set token to JwtAPI instance for subsequent calls through same instance.
func (j *JwtAPI) RequestTokenByRefreshToken(rtoken string) (Token, error) {
	err := j.requestTokenByRefreshToken(rtoken)
	if err != nil {
		return Token{}, err
	}

	j.logDebug("RequestTokenByRefreshToken", "API Object's refresh-token is (%s)", j.GetToken().RefreshToken)
	return j.GetToken(), nil
}

// Get - make HTTP GET request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *JwtAPI) Get(apipath string) (APIResult, error) {
	return j.GetURL(j.GetBaseURL() + apipath)
}

// GetURL - call given apiurl with GET method, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) GetURL(apiurl string) (APIResult, error) {
	var res APIResult
	resp, err := j.makeRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	if j.StructuredResponse {
		return getAPIResultJWT(resp)
	}
	return getRawResultJWT(resp)
}

// Post - make HTTP POST request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *JwtAPI) Post(apipath string, postdataJSON []byte) (APIResult, error) {
	return j.PostURL(j.GetBaseURL()+apipath, postdataJSON)
}

// PostURL - call given apiurl with POST method and pass data, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) PostURL(apiurl string, postdataJSON []byte) (APIResult, error) {
	var res APIResult
	if postdataJSON == nil {
		return res, fmt.Errorf("postdata is nil")
	}

	resp, err := j.makeRequest(http.MethodPost, apiurl, bytes.NewBuffer(postdataJSON))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	if j.StructuredResponse {
		return getAPIResultJWT(resp)
	}
	return getRawResultJWT(resp)
}

// Put - make HTTP PUT request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *JwtAPI) Put(apipath string, putdataJSON []byte) (APIResult, error) {
	return j.PutURL(j.GetBaseURL()+apipath, putdataJSON)
}

// PutURL - call given apiurl with PUT method and pass data, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) PutURL(apiurl string, putdataJSON []byte) (APIResult, error) {
	var res APIResult
	if putdataJSON == nil {
		return res, fmt.Errorf("putdata is nil")
	}

	resp, err := j.makeRequest(http.MethodPut, apiurl, bytes.NewBuffer(putdataJSON))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	if j.StructuredResponse {
		return getAPIResultJWT(resp)
	}
	return getRawResultJWT(resp)
}

// Delete - make HTTP DELETE request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *JwtAPI) Delete(apipath string) (APIResult, error) {
	return j.DeleteURL(j.GetBaseURL() + apipath)
}

// DeleteURL - call given apiurl with DELETE method, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) DeleteURL(apiurl string) (APIResult, error) {
	var res APIResult
	resp, err := j.makeRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	if j.StructuredResponse {
		return getAPIResultJWT(resp)
	}
	return getRawResultJWT(resp)
}

func getRawResultJWT(resp *http.Response) (APIResult, error) {
	defer resp.Body.Close()
	res := APIResult{}
	res.HTTPStatus = resp.StatusCode
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Data = err.Error()
	} else {
		res.Data = string(resbody)
	}
	res.ErrValid = false
	return res, nil
}

func getAPIResultJWT(resp *http.Response) (APIResult, error) {
	res := APIResult{}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return handleNotOK(resp), nil
	}

	json.NewDecoder(resp.Body).Decode(&res)
	res.ErrValid = true
	return res, nil
}

// func getRawResultJWT(resp *http.Response) (RawResult, error) {
// 	defer resp.Body.Close()
// 	res := RawResult{}
// 	res.HTTPStatus = resp.StatusCode
// 	resbody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		res.Data = err.Error()
// 	} else {
// 		res.Data = string(resbody)
// 	}
// 	return res, nil
// }
