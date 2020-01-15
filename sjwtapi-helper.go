package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//RequestTokenByCred call Token endpoint to get new token by passing TokenRequest data.
//It set token to JwtAPI instance for subsequent calls through same instance.
func (j *SJwtAPI) RequestTokenByCred() (Token, error) {
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
func (j *SJwtAPI) RequestTokenByRefreshToken(rtoken string) (Token, error) {
	token, err := requestTokenByRefreshToken(*j, rtoken)
	if err != nil {
		return Token{}, err
	}
	//Save token for later use
	j.token = token
	return token, nil
}

//Get - call given apiurl with GET method, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) Get(apiurl string) (APIResult, error) {
	var res APIResult
	resp, err := makeRequest(j, http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getAPIResultJWT(resp)
}

//Post - call given apiurl with POST method and pass data, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) Post(apiurl string, postdataJSON []byte) (APIResult, error) {
	var res APIResult
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

	return getAPIResultJWT(resp)
}

//Put - call given apiurl with PUT method and pass data, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) Put(apiurl string, putdataJSON []byte) (APIResult, error) {
	var res APIResult
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

	return getAPIResultJWT(resp)
}

//Delete - call given apiurl with DELETE method, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) Delete(apiurl string) (APIResult, error) {
	var res APIResult
	resp, err := makeRequest(j, http.MethodGet, apiurl, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return res, err
	}

	return getAPIResultJWT(resp)
}

//handleNotOK tried to read response from responses other than 200
// and populate ApiResult struct
func handleNotOK(resp *http.Response) APIResult {
	res := APIResult{}
	res.HTTPStatus = resp.StatusCode
	res.ErrCode = 1

	responseData, err := ioutil.ReadAll(resp.Body)
	strResp := string(responseData)

	// try to decode response to APIResult (if it has)
	if strings.Contains(strResp, "ErrText") {
		tmpres := APIResult{}
		err = jsonStringToStruct(strResp, &tmpres)
		if err == nil {
			return tmpres
		}
	}

	if err == nil {
		res.ErrText = strResp
	}
	return res
}

func getAPIResultJWT(resp *http.Response) (APIResult, error) {
	res := APIResult{}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return handleNotOK(resp), nil
	}

	json.NewDecoder(resp.Body).Decode(&res)
	return res, nil
}