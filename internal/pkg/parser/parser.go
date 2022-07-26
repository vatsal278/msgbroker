package Parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/go-playground/validator"
	"github.com/vatsal278/msgbroker/internal/model"
)

func ParseResponse(req_body io.ReadCloser, x interface{}) error {
	//Read body of the request
	body, err := ioutil.ReadAll(req_body)
	defer req_body.Close()

	if err != nil {
		log.Println(err.Error())
		return err
	}
	//Write body to struct
	err = json.Unmarshal(body, &x)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func ValidateRequest(x interface{}) error {
	validate := validator.New()
	errs := validate.Struct(x)
	if errs != nil {
		log.Print(errs.Error())
		return errs
	}
	return nil
}

func Response_Writer(status int, msg string, data interface{}) model.Response {
	var response model.Response
	response.Status = status
	response.Message = msg
	response.Data = data
	return response
}
