package parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/vatsal278/msgbroker/internal/model"
)

type Response interface {
	Update(int, string, interface{})
}

func ParseRequest(req_body io.ReadCloser, x interface{}) error {
	//Read body of the request
	body, err := ioutil.ReadAll(req_body)
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

func ValidateRequest(x interface{}) error {
	validate := validator.New()
	errs := validate.Struct(x)
	if errs != nil {
		return errs
	}
	return nil
}

func Response_Writer(w http.ResponseWriter, status int, msg string, data interface{}, r model.Response) model.Response {
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
