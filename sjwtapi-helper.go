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
	"strings"
)

//SetLogWriter Sets io.writer for logging to file
func (sj *SJwtAPI) SetLogWriter(w io.Writer) {
	multi := io.MultiWriter(w, os.Stdout)
	l := log.New(multi, "", log.LstdFlags)
	l.Println("Log output set to StdOut and writer both")
	sj.logger = l
}

//RequestTokenByCred call Token endpoint to get new token by passing TokenRequest data.
//It set token to JwtAPI instance for subsequent calls through same instance.
func (j *SJwtAPI) RequestTokenByCred() (Token, error) {
	token, err := requestTokenByLogin(*j)
	if err != nil {
		return Token{}, err
	}
	//Save token for later use
	j.token = token
	j.logDebug("RequestTokenByRefreshToken", "API Object's refresh-token is (%s)", j.GetToken().RefreshToken)
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
	j.logDebug("RequestTokenByRefreshToken", "API Object's refresh-token is (%s)", j.GetToken().RefreshToken)
	return token, nil
}

//Get - make HTTP GET request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *SJwtAPI) Get(apipath string) (APIResult, error) {
	return j.GetURL(j.GetBaseURL() + apipath)
}

//GetURL - call given apiurl with GET method, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) GetURL(apiurl string) (APIResult, error) {
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

//Post - make HTTP POST request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *SJwtAPI) Post(apipath string, postdataJSON []byte) (APIResult, error) {
	return j.PostURL(j.GetBaseURL()+apipath, postdataJSON)
}

//PostURL - call given apiurl with POST method and pass data, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) PostURL(apiurl string, postdataJSON []byte) (APIResult, error) {
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

//Put - make HTTP PUT request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *SJwtAPI) Put(apipath string, putdataJSON []byte) (APIResult, error) {
	return j.PutURL(j.GetBaseURL()+apipath, putdataJSON)
}

//PutURL - call given apiurl with PUT method and pass data, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) PutURL(apiurl string, putdataJSON []byte) (APIResult, error) {
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

//Delete - make HTTP DELETE request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (j *SJwtAPI) Delete(apipath string) (APIResult, error) {
	return j.DeleteURL(j.GetBaseURL() + apipath)
}

//DeleteURL - call given apiurl with DELETE method, auto inject Authorization Header, returns APIResult{}.
func (j *SJwtAPI) DeleteURL(apiurl string) (APIResult, error) {
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
