//Package responseWriter provides utilities to write a response to http.ResponseWriter.
package responseWriter

import (
	"encoding/json"

	"net/http"
)

//go:generate mockgen --destination=./../../mocks/mock_response.go --package=mocks github.com/vatsal278/msgbroker/pkg/responseWriter Response
type Response interface {
	Update(int, string, interface{})
}

// ResponseWriter writes a response to http.ResponseWriter in json format
func ResponseWriter(w http.ResponseWriter, status int, msg string, data interface{}, r Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	r.Update(status, msg, data)
	err := json.NewEncoder(w).Encode(r)

	return err
}
