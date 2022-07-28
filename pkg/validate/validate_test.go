package validate_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/validate"
)

var callBack = model.CallBack{
	HttpMethod:  "GET",
	CallbackUrl: "http://localhost:8083/pong",
}

func TestParser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   interface{}
		testcase      int
		expectedError interface{}
	}{
		{
			name: "SUCCESS:: validate",
			requestBody: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			testcase:      1,
			expectedError: nil,
		},
		{
			name:     "FAILURE:: validate",
			testcase: 2,
			requestBody: model.TempPublisher{
				Name:    2,
				Channel: "c4",
			},
			expectedError: errors.New("Key: 'Publisher.Name' Error:Field validation for 'Name' failed on the 'required' tag"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publisher model.Publisher
			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))

			parser.Parse(r.Body, &publisher)
			err := validate.Validate(publisher)
			if tt.testcase == 1 {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				return
			}
			//error := errors.New("Key: 'Publisher.Name' Error:Field validation for 'Name' failed on the 'required' tag")
			if err == nil {
				t.Errorf("Want: %v, Got: %v", tt.expectedError, err.Error())
			}

		})
	}
}
