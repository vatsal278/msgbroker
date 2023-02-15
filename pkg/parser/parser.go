//Package parser provides a function to parse the body of an HTTP request and write it to a struct.
package parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

//Parse helps parse the body of an HTTP request and stores it inside the struct
func Parse(reqbody io.ReadCloser, x interface{}) error {
	body, err := ioutil.ReadAll(reqbody)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &x)
	if err != nil {
		return err
	}
	return nil
}
