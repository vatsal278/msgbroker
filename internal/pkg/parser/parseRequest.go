package parser

import (
	"io"
	"log"
	"net/http"

	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
	"github.com/vatsal278/msgbroker/pkg/validate"
)

func Parse(w http.ResponseWriter, r io.ReadCloser, m interface{}) {
	err := parser.ParseRequest(r, &m)
	if err != nil {
		responseWriter.ResponseWriter(w, http.StatusInternalServerError, constants.Parse_Err, nil, &model.Response{})
		log.Println(err.Error())
		return
	}
	err = validate.Validate(&m)
	if err != nil {
		responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.Incomplete_Data, nil, &model.Response{})
		log.Println(err.Error())
		return
	}
}
