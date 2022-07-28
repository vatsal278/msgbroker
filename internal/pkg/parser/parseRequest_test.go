package parser_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/parser"
	"github.com/vatsal278/msgbroker/pkg/validate"
)

var callBack = model.CallBack{
	HttpMethod:  "GET",
	CallbackUrl: "http://localhost:8083/pong",
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		testcase         int
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: Parser",
			requestBody: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			testcase: 1,
			expectedResponse: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
		},
		{
			name:     "FAILURE:: Parser",
			testcase: 2,
			requestBody: model.Subscriber{
				CallBack: callBack,
				Channel:  "c4",
			},
			expectedResponse: model.Subscriber{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publisher model.Publisher
			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			if tt.testcase == 1 {
				err := parser.Parse(r.Body, &publisher)
				if err != nil {
					t.Error(err.Error())
				}
				err = validate.Validate(publisher)
				if err != nil {
					t.Log(err.Error())
				}
				if !reflect.DeepEqual(publisher, tt.expectedResponse) {
					t.Errorf("Want: %v, Got: %v", tt.expectedResponse, publisher)
				}
				return
			}
			err := parser.Parse(r.Body, &publisher)

			if err != nil {
				t.Log(err.Error())
			}
			err = validate.Validate(publisher)
			if err != nil {
				t.Log(err.Error())
			}
		})
	}
}
