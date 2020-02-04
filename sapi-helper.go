package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//Get - make HTTP GET request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *SAPI) Get(apipath string) (APIResult, error) {
	return a.GetURL(a.GetBaseURL() + apipath)
}

//GetURL - make HTTP GET request to given url and return APIResult{}.
func (a *SAPI) GetURL(apiurl string) (APIResult, error) {
	var res APIResult
	r, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return res, err
	}
	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getAPIResult(resp)
}

//Post - make HTTP POST request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *SAPI) Post(apipath string, postdataJSON []byte) (APIResult, error) {
	return a.PostURL(a.GetBaseURL()+apipath, postdataJSON)
}

//PostURL - make HTTP POST request to given url and post JSON data and return APIResult{}.
func (a *SAPI) PostURL(apiurl string, postdataJSON []byte) (APIResult, error) {
	var res APIResult
	if postdataJSON == nil {
		a.logMsg("APIPost", "postdata is nil")
		return res, fmt.Errorf("postdata is nil")
	}

	r, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer(postdataJSON))
	if err != nil {
		return res, err
	}
	r.Header.Set("Content-Type", "application/json")

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getAPIResult(resp)
}

//Put - make HTTP PUT request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *SAPI) Put(apipath string, putdataJSON []byte) (APIResult, error) {
	return a.PutURL(a.GetBaseURL()+apipath, putdataJSON)
}

//PutURL - make HTTP PUT request to given url and post JSON data and return APIResult{}.
func (a *SAPI) PutURL(apiurl string, putdataJSON []byte) (APIResult, error) {
	var res APIResult
	if putdataJSON == nil {
		a.logMsg("APIPut", "putdata is nil")
		return res, fmt.Errorf("putdata is nil")
	}

	r, err := http.NewRequest(http.MethodPut, apiurl, bytes.NewBuffer(putdataJSON))
	if err != nil {
		return res, err
	}
	r.Header.Set("Content-Type", "application/json")

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getAPIResult(resp)
}

//Delete - make HTTP DELETE request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *SAPI) Delete(apipath string) (APIResult, error) {
	return a.DeleteURL(a.GetBaseURL() + apipath)
}

//DeleteURL - make HTTP DELETE request to given url and return APIResult{}.
func (a *SAPI) DeleteURL(apiurl string) (APIResult, error) {
	var res APIResult
	r, err := http.NewRequest(http.MethodDelete, apiurl, nil)
	if err != nil {
		return res, err
	}
	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getAPIResult(resp)
}

func getAPIResult(resp *http.Response) (APIResult, error) {
	defer resp.Body.Close()
	res := APIResult{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res, nil
}
