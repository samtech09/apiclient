package apiclient

import (
	"bytes"
	"encoding/json"
)

//APIResult is structured response returned from APIs developed by samtech09
type APIResult struct {
	HTTPStatus int
	ErrCode    int
	ErrText    string
	//Data hold resultant JSON as string
	Data string //interface{}
}

//RawResult is unstructured response returned from any API, it could be JSON or String or other
type RawResult struct {
	HTTPStatus int
	//Response is raw response received as result of calling API
	Data string
}

//toJSON convert given interface to JSON
func toJSON(d interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

//jsonStringToStruct convert JSON string to given struct
func jsonStringToStruct(jsonstring string, dest interface{}) error {
	if err := json.Unmarshal([]byte(jsonstring), dest); err != nil {
		return err
	}
	return nil
}
