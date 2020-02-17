package apiclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// //Get - make HTTP GET request to given url and return RawResult{}.
// func (a *API) Get(apiurl string) (RawResult, error) {
// 	r, err := http.NewRequest(http.MethodGet, apiurl, nil)
// 	if err != nil {
// 		return RawResult{}, err
// 	}
// 	client := a.getClient()
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		return RawResult{}, err
// 	}
// 	return getRawResult(resp), nil
// }

//APIGet - make HTTP GET request to given url and return RawResult{}.
func (a *API) APIGet(apiurl string) (RawResult, error) {
	var res RawResult
	r, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return res, err
	}
	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getRawResult(resp), nil
}

//APIPost - make HTTP POST request to given url and post JSON data and return RawResult{}.
func (a *API) APIPost(apiurl string, postdataJSON []byte) (RawResult, error) {
	var res RawResult
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

	return getRawResult(resp), nil
}

//APIPostForm - make HTTP POST request to given url with content-type: application/x-www-form-urlencoded
func (a *API) APIPostForm(apiurl string, data map[string]string) (RawResult, error) {
	var res RawResult

	udata := url.Values{}
	for k, v := range data {
		udata.Set(k, v)
	}

	r, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBufferString(udata.Encode()))
	if err != nil {
		return res, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getRawResult(resp), nil
}

//APIPut - make HTTP PUT request to given url and post JSON data and return RawResult{}.
func (a *API) APIPut(apiurl string, putdataJSON []byte) (RawResult, error) {
	var res RawResult
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

	return getRawResult(resp), nil
}

//APIDelete - make HTTP DELETE request to given url and return RawResult{}.
func (a *API) APIDelete(apiurl string) (RawResult, error) {
	var res RawResult
	r, err := http.NewRequest(http.MethodDelete, apiurl, nil)
	if err != nil {
		return res, err
	}
	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	return getRawResult(resp), nil
}

func getRawResult(resp *http.Response) RawResult {
	defer resp.Body.Close()
	res := RawResult{}
	res.HTTPStatus = resp.StatusCode
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Data = err.Error()
	} else {
		res.Data = string(resbody)
	}
	return res
}
