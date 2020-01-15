package apiclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

//RequestTokenByCred call Token endpoint to get new token by passing TokenRequest data.
//It set token to JwtAPI instance for subsequent calls through same instance.
func (j *JwtAPI) RequestTokenByCred() (Token, error) {
	token, err := requestTokenByLogin(*j)
	if err != nil {
		return Token{}, err
	}
	//Save token for later use
	j.token = token
	return token, nil
}

//RequestTokenByRefreshToken call Token endpoint to get new token by passing existing refresh-token.
//It set token to JwtAPI instance for subsequent calls through same instance.
func (j *JwtAPI) RequestTokenByRefreshToken(rtoken string) (Token, error) {
	token, err := requestTokenByRefreshToken(*j, rtoken)
	if err != nil {
		return Token{}, err
	}
	//Save token for later use
	j.token = token
	return token, nil
}

//Get - call given apiurl with GET method, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) Get(apiurl string) (RawResult, error) {
	var res RawResult
	resp, err := makeRequest(j, http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getRawResultJWT(resp)
}

//Post - call given apiurl with POST method and pass data, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) Post(apiurl string, postdataJSON []byte) (RawResult, error) {
	var res RawResult
	if postdataJSON == nil {
		return res, fmt.Errorf("postdata is nil")
	}

	resp, err := makeRequest(j, http.MethodPost, apiurl, bytes.NewBuffer(postdataJSON))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getRawResultJWT(resp)
}

//Put - call given apiurl with PUT method and pass data, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) Put(apiurl string, putdataJSON []byte) (RawResult, error) {
	var res RawResult
	if putdataJSON == nil {
		return res, fmt.Errorf("putdata is nil")
	}

	resp, err := makeRequest(j, http.MethodPut, apiurl, bytes.NewBuffer(putdataJSON))
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getRawResultJWT(resp)
}

//Delete - call given apiurl with DELETE method, auto inject Authorization Header, returns RawResult{}.
func (j *JwtAPI) Delete(apiurl string) (RawResult, error) {
	var res RawResult
	resp, err := makeRequest(j, http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getRawResultJWT(resp)
}

func getRawResultJWT(resp *http.Response) (RawResult, error) {
	defer resp.Body.Close()
	res := RawResult{}
	res.HTTPStatus = resp.StatusCode
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Data = err.Error()
	} else {
		res.Data = string(resbody)
	}
	return res, nil
}
