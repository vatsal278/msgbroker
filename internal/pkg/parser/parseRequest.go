// Package parseRequest provides functions that helps parse the body of an HTTP request, write it to a struct and validate all the fields
package parseRequest

import (
	"io"

	"github.com/vatsal278/msgbroker/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/validate"
)

// ParseAndValidateRequest helps parse the body of an HTTP request, store it in a struct and validate the fields
func ParseAndValidateRequest(r io.ReadCloser, m interface{}) error {
	err := parser.Parse(r, &m)
	if err != nil {
		return err
	}
	err = validate.Validate(m)
	if err != nil {
		return err
	}
	return nil
}
