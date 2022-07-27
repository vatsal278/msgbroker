package responseWriter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vatsal278/msgbroker/internal/model"
)

type Response interface {
	Update(int, string, interface{})
}

func ResponseWriter(w http.ResponseWriter, status int, msg string, data interface{}, r model.Response) model.Response {
	var response model.Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	r.Update(status, msg, data)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Print(err.Error())
	}
	return response
}
