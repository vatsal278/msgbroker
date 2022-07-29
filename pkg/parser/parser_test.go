package parser_test

import (
	"bytes"
	"encoding/json"
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
		testcase         int
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			testcase: 1,
			expectedResponse: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
		},
		{
			name:     "FAILURE:: Parser",
			testcase: 2,
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			expectedResponse: testStructFail{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var teststruct testStruct
			var teststructfail testStructFail
			//var subscriber model.Subscriber
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			bytes.NewReader(jsonValue)
			if tt.testcase == 1 {
				err := parser.Parse(r.Body, &teststruct)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(teststruct, tt.expectedResponse) {
					t.Errorf("Want: %v, Got: %v", tt.expectedResponse, &teststruct)
				}
				return
			}
			err := parser.Parse(r.Body, &teststructfail)
			if err != nil {
				t.Log(err.Error())
			} else {
				t.Errorf("Want: %v, Got: %v", "error", nil)
			}
			if !reflect.DeepEqual(teststructfail, tt.expectedResponse) {
				t.Errorf("Want: %v, Got: %v", tt.expectedResponse, teststructfail)
			}
		})
	}
}
