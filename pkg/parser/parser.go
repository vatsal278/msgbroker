package parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//Parse helps parse the body of an HTTP request and write it to a struct
func Parse(reqbody io.ReadCloser, x interface{}) error {
	//Read body of the request
	body, err := ioutil.ReadAll(reqbody)
	//defer req_body.Close()

	if err != nil {
		return err
	}
	//Write body to struct
	err = json.Unmarshal(body, &x)
	if err != nil {
		return err
	}
	return nil
}
