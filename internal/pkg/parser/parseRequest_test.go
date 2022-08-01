package parseRequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

type testStruct struct {
	Name    string
	Channel string
}
type testStructFail struct {
	Name    int
	Channel int
}

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		setupFunc        func(r *http.Request, publisher model.Publisher)
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: ParseRequest",
			requestBody: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupFunc: func(r *http.Request, publisher model.Publisher) {
				var expectedResponse = model.Publisher{
					Name:    "publisher1",
					Channel: "c4",
				}
				err := parser.Parse(r.Body, &publisher)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				err = validate.Validate(publisher)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(publisher, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, publisher)
				}
			},
			expectedResponse: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
		},
		{
			name: "FAILURE:: Parse",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupFunc: func(r *http.Request, publisher model.Publisher) {
				var teststructfail testStructFail
				err := parser.Parse(r.Body, &teststructfail)
				expectedResponse := testStructFail{}
				if err != nil {
					t.Log(err.Error())
				} else {
					t.Errorf("Want: %v, Got: %v", "error", nil)
				}
				if !reflect.DeepEqual(teststructfail, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, teststructfail)
				}
			},
		},
		{
			name: "FAILURE:: validate function",
			requestBody: model.TempPublisher{
				Name:    2,
				Channel: "c4",
			},
			setupFunc: func(r *http.Request, publisher model.Publisher) {
				err := parser.Parse(r.Body, &publisher)
				if err != nil {
					t.Log(err.Error())
				}
				err = validate.Validate(publisher)

				if err != nil {
					t.Log(err.Error())
				} else {
					t.Errorf("Want: %v, Got: %v", errors.New("Key: 'Publisher.Name' Error:Field validation for 'Name' failed on the 'required' tag"), err.Error())
				}
			},
			expectedResponse: model.Publisher{
				Channel: "c4",
			},
		},
	}
	//parse succes, parse succes validate failure, parse failure
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publisher model.Publisher
			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			tt.setupFunc(r, publisher)
		})
	}
}
