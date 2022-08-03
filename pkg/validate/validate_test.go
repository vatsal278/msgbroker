package validate

import (
	"errors"
	"testing"
)

type Publisher struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

func TestParser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   interface{}
		validateFunc  func(error)
		expectedError interface{}
	}{
		{
			name: "SUCCESS:: validate",
			requestBody: Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			validateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name: "FAILURE:: validate:: Incomplete Input Details",
			requestBody: Publisher{

				Channel: "c4",
			},
			validateFunc: func(err error) {
				if err != nil {
					return
				} else {
					t.Errorf("Want: %v, Got: %v", errors.New("Key: 'Publisher.Name' Error:Field validation for 'Name' failed on the 'required' tag"), err)
				}
			},
		},
	}
	//remove the http request logic and directly pass in local struct to validate
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var publisher = tt.requestBody
			err := Validate(publisher)
			tt.validateFunc(err)
		})
	}
}
