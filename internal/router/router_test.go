package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/msgbroker/internal/constants"
	"github.com/vatsal278/msgbroker/internal/model"
)

type temp_struct struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func TestRegisterPublisher(t *testing.T) {
	var publisher = model.Publisher{
		Name:    "publisher1",
		Channel: "c4",
	}
	var dummy = model.TempPublisher{
		Name:    1,
		Channel: "c4",
	}
	tests := []struct {
		name              string
		requestBody       interface{}
		ErrorCase         bool
		expected_response temp_struct
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: publisher,
			expected_response: temp_struct{
				Status:  http.StatusCreated,
				Message: "Successfully Registered as publisher to the channel",
				Data:    nil,
			},
		},
		{
			name:        "FAILURE:: Register Publisher",
			requestBody: dummy,
			expected_response: temp_struct{
				Status:  http.StatusBadRequest,
				Message: constants.IncompleteData,
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			t.Log(w.Code)
			jsonValue, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))
			router := Router()
			router.ServeHTTP(w, r)
			response_body, error := ioutil.ReadAll(w.Body)
			if error != nil {
				t.Error(error.Error())
			}
			var response temp_struct
			err := json.Unmarshal(response_body, &response)
			if err != nil {
				t.Error(error.Error())
			}
			if !reflect.DeepEqual(response, tt.expected_response) {
				t.Errorf("Want: %v, Got: %v", tt.expected_response, response)
			}
		})
	}
}
