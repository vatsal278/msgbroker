package parser_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/parser"
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

func TestParser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		setupFunc        func(r *http.Request)
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupFunc: func(r *http.Request) {
				var teststruct testStruct
				err := parser.Parse(r.Body, &teststruct)
				expectedResponse := testStruct{
					Name:    "publisher1",
					Channel: "c4",
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(teststruct, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, &teststruct)
				}
			},
		},
		{
			name: "FAILURE:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupFunc: func(r *http.Request) {
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
			expectedResponse: testStructFail{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			//bytes.NewReader(jsonValue)
			tt.setupFunc(r)

		})
	}
}
