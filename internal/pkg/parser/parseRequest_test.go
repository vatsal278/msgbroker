package parseRequest

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParseAndValidateRequest(t *testing.T) {
	type testStruct struct {
		Field1 int `json:"field_1" validate:"required"`
	}
	tests := []struct {
		name        string
		requestBody io.ReadCloser
		validation  func(error, testStruct)
	}{
		{
			name:        "FAILURE:: Parse failure::json unmarshall error",
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
		{
			name:        "FAILURE:: Parse success::validation failure",
			requestBody: io.NopCloser(bytes.NewBuffer([]byte("{}"))),
			validation: func(err error, x testStruct) {
				if !strings.Contains(err.Error(), "failed on the 'required' tag") {
					t.Errorf("Want: %v, Got: %v", "failed on the 'required' tag", err.Error())
				}
				expectedResponse := testStruct{}
				if !reflect.DeepEqual(x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, x)
				}
			},
		},
		{
			name:        "SUCCESS",
			requestBody: io.NopCloser(bytes.NewBuffer([]byte(`{"field_1": 12}`))),
			validation: func(err error, x testStruct) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				expectedResponse := testStruct{
					Field1: 12,
				}
				if !reflect.DeepEqual(x, expectedResponse) {
					t.Errorf("Want: %v, Got: %v", expectedResponse, x)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := testStruct{}
			err := ParseAndValidateRequest(tt.requestBody, &temp)
			tt.validation(err, temp)
		})
	}
}
