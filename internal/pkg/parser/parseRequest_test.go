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
)

type Publisher struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}
type testStruct struct {
	Name    string
	Channel string
}
type testStructFail struct {
	Name    int
	Channel int
}

func TestParseAndValidateRequest(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		setupFunc        func(r *http.Request, publisher Publisher)
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: ParseRequest",
			requestBody: model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupFunc: func(r *http.Request, publisher Publisher) {
				var expectedResponse = Publisher{
					Name:    "publisher1",
					Channel: "c4",
				}
				err := ParseAndValidateRequest(r.Body, publisher)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(publisher, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, &publisher)
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
			setupFunc: func(r *http.Request, publisher Publisher) {
				var teststructfail testStructFail
				err := ParseAndValidateRequest(r.Body, teststructfail)
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
			requestBody: testStructFail{
				Name: 2,
			},
			setupFunc: func(r *http.Request, publisher Publisher) {
				err := ParseAndValidateRequest(r.Body, publisher)
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
			var publisher Publisher
			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			tt.setupFunc(r, publisher)
		})
	}
}
