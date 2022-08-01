package validate

import (
	"errors"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
)

var callBack = model.CallBack{
	HttpMethod:  "GET",
	CallbackUrl: "http://localhost:8083/pong",
}

type Publisher struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}
type tempPublisher struct {
	Name    interface{} `validate:"required"`
	Channel string      `form:"channel" json:"channel" validate:"required"`
}

func TestParser(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   interface{}
		setupFunc     func()
		expectedError interface{}
	}{
		{
			name: "SUCCESS:: validate",
			setupFunc: func() {
				var publisher = Publisher{
					Name:    "publisher1",
					Channel: "c4",
				}
				err := Validate(publisher)
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name: "FAILURE:: validate",
			setupFunc: func() {
				var publisher = tempPublisher{
					Channel: "c4",
				}
				err := Validate(publisher)
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
			tt.setupFunc()
		})
	}
}
