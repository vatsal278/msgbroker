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
		setupModel       func() interface{}
		validation       func(error, interface{})
		expectedResponse interface{}
	}{
		{
			name: "SUCCESS:: Parser",
			requestBody: testStruct{
				Name:    "publisher1",
				Channel: "c4",
			},
			setupModel: func() interface{} {
				var teststruct testStruct
				return teststruct
			},
			validation: func(err error, x interface{}) {
				t.Log(x.(testStruct))
				expectedResponse := testStruct{
					Name:    "publisher1",
					Channel: "c4",
				}
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if !reflect.DeepEqual(&x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, &x)
				}
			},
		},
		/*{
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
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			temp := tt.setupModel()
			x := temp.(testStruct)

			err := Parse(r.Body, x)

			tt.validation(err, x)

		})
	}
}
