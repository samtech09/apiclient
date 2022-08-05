package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// //Get - make HTTP GET request to given url and return RawResult{}.
// func (a *API) Get(apiurl string) (APIResult error) {
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

//SetLogWriter Sets io.writer for logging to file
func (a *API) SetLogWriter(w io.Writer) {
	multi := io.MultiWriter(w, os.Stdout)
	l := log.New(multi, "", log.LstdFlags)
	l.Println("Log output set to StdOut and writer both")
	a.logger = l
}

//Get - make HTTP GET request to given api path and return RawResult{}.
func (a *API) Get(apipath string) (APIResult, error) {
	return a.GetURL(a.GetBaseURL() + apipath)
}

//GetURL - make HTTP GET request to given url and return RawResult{}.
func (a *API) GetURL(apiurl string) (APIResult, error) {
	var res APIResult
	r, err := http.NewRequest(http.MethodGet, apiurl, nil)
	if err != nil {
		return res, err
	}
	a.injectHeaders(r)
	if a.UseBasicAuth {
		r.SetBasicAuth(a.BasicAuthUser, a.BasicAuthPwd)
	}

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	if a.StructuredResponse {
		return getAPIResult(resp)
	}
	return getRawResult(resp), nil
}

//Post - make HTTP POST request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *API) Post(apipath string, postdataJSON []byte) (APIResult, error) {
	return a.PostURL(a.GetBaseURL()+apipath, postdataJSON)
}

//PostURL - make HTTP POST request to given url and post JSON data and return RawResult{}.
func (a *API) PostURL(apiurl string, postdataJSON []byte) (APIResult, error) {
	var res APIResult
	if postdataJSON == nil {
		a.logMsg("APIPost", "postdata is nil")
		return res, fmt.Errorf("postdata is nil")
	}

	r, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBuffer(postdataJSON))
	if err != nil {
		return res, err
	}
	a.injectHeaders(r)
	r.Header.Set("Content-Type", "application/json")
	if a.UseBasicAuth {
		r.SetBasicAuth(a.BasicAuthUser, a.BasicAuthPwd)
	}

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	if a.StructuredResponse {
		return getAPIResult(resp)
	}
	return getRawResult(resp), nil
}

//PostForm - make HTTP POST request to given url with content-type: application/x-www-form-urlencoded
func (a *API) PostForm(apiurl string, data map[string]string) (APIResult, error) {
	var res APIResult

	udata := url.Values{}
	for k, v := range data {
		udata.Set(k, v)
	}

	r, err := http.NewRequest(http.MethodPost, apiurl, bytes.NewBufferString(udata.Encode()))
	if err != nil {
		return res, err
	}
	a.injectHeaders(r)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if a.UseBasicAuth {
		r.SetBasicAuth(a.BasicAuthUser, a.BasicAuthPwd)
	}

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	if a.StructuredResponse {
		return getAPIResult(resp)
	}
	return getRawResult(resp), nil
}

//Put - make HTTP PUT request to given api path, post JSON data and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *API) Put(apipath string, putdataJSON []byte) (APIResult, error) {
	return a.PutURL(a.GetBaseURL()+apipath, putdataJSON)
}

//PutURL - make HTTP PUT request to given url and post JSON data and return RawResult{}.
func (a *API) PutURL(apiurl string, putdataJSON []byte) (APIResult, error) {
	var res APIResult
	if putdataJSON == nil {
		a.logMsg("APIPut", "putdata is nil")
		return res, fmt.Errorf("putdata is nil")
	}

	r, err := http.NewRequest(http.MethodPut, apiurl, bytes.NewBuffer(putdataJSON))
	if err != nil {
		return res, err
	}
	a.injectHeaders(r)
	r.Header.Set("Content-Type", "application/json")
	if a.UseBasicAuth {
		r.SetBasicAuth(a.BasicAuthUser, a.BasicAuthPwd)
	}

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	if a.StructuredResponse {
		return getAPIResult(resp)
	}
	return getRawResult(resp), nil
}

//Delete - make HTTP DELETE request to given api path and return APIResult{}. ResourceAPIBaseURL will be prepended.
func (a *API) Delete(apipath string) (APIResult, error) {
	return a.DeleteURL(a.GetBaseURL() + apipath)
}

//DeleteURL - make HTTP DELETE request to given url and return RawResult{}.
func (a *API) DeleteURL(apiurl string) (APIResult, error) {
	var res APIResult
	r, err := http.NewRequest(http.MethodDelete, apiurl, nil)
	if err != nil {
		return res, err
	}
	a.injectHeaders(r)
	if a.UseBasicAuth {
		r.SetBasicAuth(a.BasicAuthUser, a.BasicAuthPwd)
	}

	client := a.getClient()
	resp, err := client.Do(r)
	if err != nil {
		return res, err
	}

	if a.StructuredResponse {
		return getAPIResult(resp)
	}
	return getRawResult(resp), nil
}

func (a *API) injectHeaders(r *http.Request) {
	if len(a.headers) > 0 {
		for k, v := range a.headers {
			r.Header.Set(k, v)
		}
	}
}

func getRawResult(resp *http.Response) APIResult {
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
	return res
}

func getAPIResult(resp *http.Response) (APIResult, error) {
	defer resp.Body.Close()
	res := APIResult{}
	json.NewDecoder(resp.Body).Decode(&res)
	res.ErrValid = true
	return res, nil
}
