package apiclient

import (
	"net/http"
	"time"
)

type iapi interface {
	InsecureSSLEnabled() bool
	DebugEnabled() bool
	GetTimeout() time.Duration
	GetBaseURL() string

	getClient() *http.Client
	logMsg(methodname, format string, msg ...interface{})
}

//API - provide functions to call APIs using GET, POST, PUT, DELETE methods to any API. It will return RawReuslt{}.
//API response will be read and set as string into RawResult.Data that client can parse.
type API struct {
	AllowInsecureSSL   bool
	Timeout            time.Duration
	Debug              bool
	ResourceAPIBaseURL string
}

//SAPI - allow to make calls to Structured APIs using GET, POST, PUT, DELETE methods which itself return response as APIReuslt{}
type SAPI struct {
	AllowInsecureSSL   bool
	Timeout            time.Duration
	Debug              bool
	ResourceAPIBaseURL string
}

func (j API) InsecureSSLEnabled() bool {
	return j.AllowInsecureSSL
}
func (j API) DebugEnabled() bool {
	return j.Debug
}
func (j API) GetTimeout() time.Duration {
	return j.Timeout
}
func (j API) getClient() *http.Client {
	return getClient(j.InsecureSSLEnabled(), j.GetTimeout())
}
func (j API) GetBaseURL() string {
	return j.ResourceAPIBaseURL
}
func (j API) logMsg(methodname, format string, msg ...interface{}) {
	logMsg(j.DebugEnabled(), methodname, format, msg...)
}

func (j SAPI) InsecureSSLEnabled() bool {
	return j.AllowInsecureSSL
}
func (j SAPI) DebugEnabled() bool {
	return j.Debug
}
func (j SAPI) GetTimeout() time.Duration {
	return j.Timeout
}
func (j SAPI) GetBaseURL() string {
	return j.ResourceAPIBaseURL
}
func (j SAPI) getClient() *http.Client {
	return getClient(j.InsecureSSLEnabled(), j.GetTimeout())
}
func (j SAPI) logMsg(methodname, format string, msg ...interface{}) {
	logMsg(j.DebugEnabled(), methodname, format, msg...)
}
