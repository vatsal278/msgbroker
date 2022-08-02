package parser

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testStruct struct {
	Name    string
	Channel string
}
type testStructFail struct {
	Name    int
	Channel int
}

//var teststruct testStruct

func TestParser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      interface{}
		setupModel       interface{}
		validation       func(error, interface{})
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupModel: testStruct{},
			validation: func(err error, x interface{}) {
				expectedResponse := testStruct{
					Name:    "publisher1",
					Channel: "c4",
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, x)
				}
			},
		},
		{
			name: "FAILURE:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupModel: testStructFail{},
			validation: func(err error, x interface{}) {
				expectedResponse := testStructFail{}
				if err != nil {
					t.Log(err.Error())
				} else {
					t.Errorf("Want: %v, Got: %v", "error", nil)
				}
				if !reflect.DeepEqual(&x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, &x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			var temp = tt.setupModel
			err := Parse(r.Body, &temp)
			tt.validation(err, &temp)
		})
	}
}
