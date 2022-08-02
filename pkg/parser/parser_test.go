package parser

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

type testStruct struct {
	Field1 int `json:"field_1"`
}

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		requestBody io.ReadCloser
		validation  func(error, testStruct)
	}{
		{
			name:        "SUCCESS::Parse success",
			requestBody: io.NopCloser(bytes.NewBuffer([]byte(`{"field_1": 23}`))),
			validation: func(err error, x testStruct) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				expectedResponse := testStruct{
					Field1: 23,
				}
				if !reflect.DeepEqual(x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, x)
				}
			},
		},
		{
			name:        "FAILURE:: json unmarshall failure",
			requestBody: io.NopCloser(bytes.NewBuffer([]byte(`{"field_1": "23"}`))),
			validation: func(err error, x testStruct) {
				if !strings.Contains(err.Error(), "cannot unmarshal string into Go struct field") {
					t.Errorf("Want: %v, Got: %v", "cannot unmarshal string into Go struct field", err.Error())
				}
				expectedResponse := testStruct{}
				if !reflect.DeepEqual(x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, x)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := testStruct{}
			err := Parse(tt.requestBody, &temp)
			tt.validation(err, temp)
		})
	}
}
